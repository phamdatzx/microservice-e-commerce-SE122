"""
E-commerce AI chatbot agent — RAG + Function Calling with GPT-4o.

Architecture
~~~~~~~~~~~~

    User query
        │
        ▼
    ┌─────────────────────────────────────────────┐
    │  LangChain Agent  (GPT-4o, temperature=0)   │
    │                                             │
    │  Tools available:                           │
    │   • search_products   (RAG → Qdrant)        │
    │   • get_product_by_id (backend API)         │
    │   • get_product_reviews                     │
    │   • get_my_orders                           │
    │   • get_seller_info_by_id                   │
    │   • get_my_vouchers                         │
    │                                             │
    │  The LLM decides which tools to call,       │
    │  interprets the results, and generates a    │
    │  natural-language answer.                   │
    └─────────────────────────────────────────────┘

Public API
----------
- ``create_agent_executor()`` → ``AgentExecutor`` (singleton-friendly)
- ``chat(query)`` → ``str``   (one-shot convenience)
"""

from __future__ import annotations

import json
import logging
from functools import lru_cache

from langchain_openai import ChatOpenAI
from langchain_core.prompts import ChatPromptTemplate, MessagesPlaceholder
from langchain_core.tools import tool
from langchain_classic.agents import AgentExecutor, create_tool_calling_agent

from app.core.config import get_settings
from app.services.rag_service import retrieve_products_filtered
from app.tools.backend_tools import ALL_TOOLS as BACKEND_TOOLS

logger = logging.getLogger(__name__)

# ─── System prompt ───────────────────────────────────────────────────────

SYSTEM_PROMPT = """\
Bạn là trợ lý mua sắm AI của một sàn thương mại điện tử.

NHIỆM VỤ:
• Giúp khách hàng tìm kiếm sản phẩm, so sánh giá, xem đánh giá,
  theo dõi đơn hàng, xem giỏ hàng, và giải đáp mọi thắc mắc về mua sắm.
• Trả lời bằng tiếng Việt, thân thiện, ngắn gọn và chính xác.

THÔNG TIN NGƯỜI DÙNG:
Mỗi tin nhắn có thể bắt đầu bằng một block "=== THÔNG TIN NGƯỜI DÙNG ===" do
hệ thống cung cấp. Thông tin này đã được xác thực — hãy tin tưởng và sử dụng
trực tiếp khi gọi tool. KHÔNG bao giờ hỏi lại người dùng về user_id.

HƯỚNG DẪN DÙNG TOOL:
1. Tìm kiếm sản phẩm (semantic):
   → `search_products` — tìm theo ngữ nghĩa, có hỗ trợ lọc giá, rating, kho.
     `query` chỉ chứa từ khóa sản phẩm (ví dụ: "iPhone", "máy giặt").
2. Tìm kiếm có bộ lọc / sắp xếp:
   → `search_product_catalog` — lọc giá, danh mục, sắp xếp theo bán chạy/rating.
3. Chi tiết sản phẩm (đã biết ID):
   → `get_product_by_id`
4. Đánh giá / review sản phẩm:
   → `get_product_reviews`
5. Thông tin về một shop/người bán:
   a. Tổng quan về shop (shop bán gì, có những loại hàng nào):
      → `get_seller_categories` — ƯU TIÊN dùng trước, trả về danh mục hàng
        mà người bán tự tạo, phản ánh trực tiếp cách họ tổ chức cửa hàng.
   b. Danh sách sản phẩm cụ thể của shop (khi người dùng muốn xem sản phẩm):
      → `get_products_by_seller` — chỉ gọi sau khi đã có tổng quan từ 5a,
        hoặc khi người dùng yêu cầu xem sản phẩm rõ ràng.
6. Voucher / mã giảm giá của shop:
   → `get_seller_vouchers`
7. Danh mục sản phẩm (nền tảng):
   → `get_categories` (dùng để lấy category_id cho bộ lọc)
8. Đơn hàng của người dùng:
   → `get_my_orders` (bắt buộc truyền user_id từ context)
9. Giỏ hàng:
   → `get_my_cart` hoặc `get_cart_count` (bắt buộc truyền user_id từ context)
10. Kiểm tra đã mua hàng chưa:
    → `verify_purchase`
11. Thông tin biến thể sản phẩm (từ ID biến thể):
    → `get_variants_by_ids`

QUY TẮC:
• CHỈ trả lời dựa trên dữ liệu từ các tool. KHÔNG bịa thông tin.
• Nếu không tìm thấy kết quả, hãy nói rõ: "Tôi không tìm thấy …".
• Hiển thị giá bằng đơn vị VNĐ (ví dụ: 1,500,000đ).
• Khi liệt kê sản phẩm, hiển thị: tên, giá, rating, tình trạng kho.
• Với mọi tool yêu cầu user_id, luôn dùng giá trị từ block thông tin người dùng.
"""


# ─── RAG as a tool (with filtering) ─────────────────────────────────────

@tool
def search_products(
    query: str,
    price_max: int | None = None,
    price_min: int | None = None,
    min_rating: float | None = None,
    in_stock: bool = True,
) -> str:
    """Search for products in the e-commerce catalog using semantic search
    with optional payload filters.

    IMPORTANT: `query` should contain ONLY product keywords
    (e.g. "iPhone", "máy giặt", "laptop gaming").
    Do NOT put price or other numeric constraints in the query string —
    use the dedicated filter parameters instead.

    Args:
        query: Product keywords only (e.g. "iPhone", "máy giặt").
        price_max: Maximum price in VND. E.g. user says "dưới 10 triệu" → price_max=10000000.
        price_min: Minimum price in VND. E.g. user says "trên 5 triệu" → price_min=5000000.
        min_rating: Minimum star rating (1-5). E.g. user says "đánh giá tốt" → min_rating=4.0.
        in_stock: Only return in-stock products. Default True.
    """
    settings = get_settings()
    docs = retrieve_products_filtered(
        query,
        k=settings.RAG_TOP_K,
        price_max=price_max,
        price_min=price_min,
        in_stock=in_stock,
        min_rating=min_rating,
    )

    if not docs:
        return "Không tìm thấy sản phẩm nào phù hợp với tiêu chí tìm kiếm."

    results = []
    for i, doc in enumerate(docs, 1):
        meta = doc.metadata
        results.append(
            f"[{i}] {meta.get('name', 'N/A')}\n"
            f"    ID: {meta.get('product_id', 'N/A')}\n"
            f"    Giá: {meta.get('price_min', '?')} – {meta.get('price_max', '?')}\n"
            f"    Rating: {meta.get('rating', 'N/A')} ({meta.get('rate_count', 0)} đánh giá)\n"
            f"    Tồn kho: {meta.get('stock', 0)}\n"
            f"    Đã bán: {meta.get('sold_count', 0)}\n"
            f"    Danh mục: {', '.join(meta.get('category_names', []))}\n"
            f"    Mô tả: {doc.page_content[:200]}"
        )

    return "\n\n".join(results)


# ─── Combine all tools ──────────────────────────────────────────────────

ALL_TOOLS = [search_products] + BACKEND_TOOLS


# ─── Agent factory ───────────────────────────────────────────────────────

def _build_prompt() -> ChatPromptTemplate:
    """Build the chat prompt with system message and placeholders."""
    return ChatPromptTemplate.from_messages(
        [
            ("system", SYSTEM_PROMPT),
            MessagesPlaceholder(variable_name="chat_history", optional=True),
            ("human", "{input}"),
            MessagesPlaceholder(variable_name="agent_scratchpad"),
        ]
    )


def _build_llm() -> ChatOpenAI:
    """Instantiate GPT-4o with temperature=0."""
    settings = get_settings()
    return ChatOpenAI(
        model=settings.OPENAI_MODEL,
        temperature=settings.OPENAI_TEMPERATURE,
        openai_api_key=settings.OPENAI_API_KEY,
    )


@lru_cache
def create_agent_executor() -> AgentExecutor:
    """Create and return a ready-to-use ``AgentExecutor``.

    The executor is cached so only one instance is created per process.
    """
    llm = _build_llm()
    prompt = _build_prompt()

    agent = create_tool_calling_agent(
        llm=llm,
        tools=ALL_TOOLS,
        prompt=prompt,
    )

    executor = AgentExecutor(
        agent=agent,
        tools=ALL_TOOLS,
        verbose=True,              # debugging output
        handle_parsing_errors=True,
        max_iterations=10,
    )

    logger.info(
        "Agent created — model=%s, tools=%s",
        _build_llm().model_name,
        [t.name for t in ALL_TOOLS],
    )
    return executor


# ─── Convenience API ─────────────────────────────────────────────────────

def chat(query: str, chat_history: list | None = None) -> str:
    """Send a single query to the agent and return the text response.

    Parameters
    ----------
    query : str
        The user's message.
    chat_history : list, optional
        List of previous ``(human, ai)`` message tuples for multi-turn
        conversations.  Defaults to empty (single-turn).

    Returns
    -------
    str
        The agent's natural-language response.
    """
    executor = create_agent_executor()
    result = executor.invoke(
        {
            "input": query,
            "chat_history": chat_history or [],
        }
    )
    return result["output"]

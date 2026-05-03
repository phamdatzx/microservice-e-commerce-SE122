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
from app.services.rag_service import retrieve_products
from app.tools.backend_tools import ALL_TOOLS as BACKEND_TOOLS

logger = logging.getLogger(__name__)

# ─── System prompt ───────────────────────────────────────────────────────

SYSTEM_PROMPT = """\
Bạn là trợ lý mua sắm AI của một sàn thương mại điện tử.

NHIỆM VỤ:
• Giúp khách hàng tìm kiếm sản phẩm, so sánh giá, xem đánh giá,
  theo dõi đơn hàng, và giải đáp mọi thắc mắc liên quan đến mua sắm.
• Trả lời bằng tiếng Việt, thân thiện, ngắn gọn và chính xác.

CÁCH LÀM VIỆC:
1. Khi khách hỏi về sản phẩm → dùng tool `search_products` để tìm kiếm
   trong cơ sở dữ liệu sản phẩm trước, rồi trả lời dựa trên kết quả.
2. Khi cần thông tin chi tiết về một sản phẩm cụ thể (đã biết ID)
   → dùng `get_product_by_id`.
3. Khi cần đánh giá / review → dùng `get_product_reviews`.
4. Khi khách hỏi về đơn hàng → dùng `get_my_orders`.
5. Khi khách hỏi về người bán → dùng `get_seller_info_by_id`.
6. Khi khách hỏi về voucher / mã giảm giá → dùng `get_my_vouchers`.

QUY TẮC:
• CHỈ trả lời dựa trên dữ liệu từ các tool. KHÔNG bịa thông tin.
• Nếu không tìm thấy kết quả, hãy nói rõ: "Tôi không tìm thấy …".
• Hiển thị giá bằng đơn vị VNĐ (ví dụ: 1,500,000đ).
• Khi liệt kê sản phẩm, hiển thị: tên, giá, rating, tình trạng kho.
"""


# ─── RAG as a tool ───────────────────────────────────────────────────────

@tool
def search_products(query: str) -> str:
    """Search for products in the e-commerce catalog using semantic search.

    Use this tool when the user asks about products, wants to find items,
    compare products, or asks general product-related questions.
    The query should be a natural-language description of what the user
    is looking for.

    Args:
        query: Natural-language search query (e.g. "máy giặt giá rẻ").
    """
    settings = get_settings()
    docs = retrieve_products(query, k=settings.RAG_TOP_K)

    if not docs:
        return "Không tìm thấy sản phẩm nào phù hợp."

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

import os
import json
from langchain_openai import ChatOpenAI
from langchain_core.prompts import PromptTemplate
from sentence_transformers import SentenceTransformer
from qdrant_client import QdrantClient
from qdrant_client.models import Filter, FieldCondition, MatchValue

# Định nghĩa trực tiếp thay vì import qua thư mục chứa dấu gạch ngang "ai-embedding"
QDRANT_HOST = os.getenv("QDRANT_HOST", "localhost")
QDRANT_PORT = int(os.getenv("QDRANT_PORT", 6333))
QDRANT_PRODUCT_COLLECTION = os.getenv("QDRANT_PRODUCT_COLLECTION", "products")
QDRANT_RATING_COLLECTION = os.getenv("QDRANT_RATING_COLLECTION", "ratings")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "BAAI/bge-m3")

OPENROUTER_BASE_URL = os.getenv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1")
OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY", "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")
# OpenRouter đã gỡ bản free của Llama 3.1 8B, nên ta dùng Llama 3.2 3B Instruct Free hoặc Google Gemini 2.0
LLM_MODEL = os.getenv("LLM_MODEL", "meta-llama/Llama-3.3-70B-Instruct")

# 1. Khởi tạo LLM trỏ qua OpenRouter (Llama 3 8B)
llm = ChatOpenAI(
    openai_api_key=OPENROUTER_API_KEY,
    openai_api_base=OPENROUTER_BASE_URL,
    model_name=LLM_MODEL,
    temperature=0.3
)

# 2. Khởi tạo Vector Search Model & Qdrant Client
print(f"Loading embedding model {EMBEDDING_MODEL} for RAG...")
embedding_model = SentenceTransformer(EMBEDDING_MODEL)
qdrant_client = QdrantClient(host=QDRANT_HOST, port=QDRANT_PORT)

# 3. Prompts
PRODUCT_PROMPT = PromptTemplate(
    template="""Bạn là trợ lý ảo hỗ trợ khách hàng của hệ thống e-commerce. Hãy trả lời câu hỏi dựa trên thông tin sản phẩm đính kèm.
Thông tin sản phẩm: {context}

Câu hỏi của khách: {query}
Câu trả lời (hãy trả lời ngắn gọn, lịch sự, đúng trọng tâm):""",
    input_variables=["context", "query"]
)

REVIEW_SUMMARY_PROMPT = PromptTemplate(
    template="""Bạn là trợ lý ảo. Nhiệm vụ của bạn là đọc các bài đánh giá của khách hàng về một sản phẩm, sau đó tóm tắt lại tổng quan khen/chê.
Các đánh giá: {context}

Câu hỏi của khách: {query}
Tóm tắt đánh giá (Chỉ dựa vào thông tin cung cấp ở trên):""",
    input_variables=["context", "query"]
)

# -------------------------------------------------------------
# Chức năng 1: Hỏi đáp về chi tiết Sản phẩm
# -------------------------------------------------------------
def ask_product_detail(query: str) -> str:
    # Bước 1: Vectorize query
    vector = embedding_model.encode(query).tolist()
    
    # Bước 2: Tìm kiếm 3 sản phẩm phù hợp nhất với Qdrant
    search_result = qdrant_client.query_points(
        collection_name=QDRANT_PRODUCT_COLLECTION,
        query=vector,
        limit=3
    )
    
    # Bước 3: Thu thập Context
    context_chunks = []
    for hit in search_result.points:
        p = hit.payload
        chunk = f"[Sản phẩm: {p.get('name', 'N/A')} - ID: {p.get('id', 'N/A')} - Giá: {p.get('price_min')} đến {p.get('price_max')} - Mô tả: {p.get('description', '')}]"
        context_chunks.append(chunk)
    
    context_text = "\n".join(context_chunks)
    
    # Bước 4: Gọi LLM
    prompt_str = PRODUCT_PROMPT.format(context=context_text, query=query)
    response = llm.invoke(prompt_str)
    
    return response.content

# -------------------------------------------------------------
# Chức năng 2: Tóm tắt Đánh giá rành cho 1 loại SP
# -------------------------------------------------------------
def ask_review_summary(query: str) -> str:
    # Thường User sẽ hỏi: "Đánh giá của máy nghe nhạc abc thế nào?"
    # Nên ta phải dùng query đó tìm Product chính xác nhất trước
    vector = embedding_model.encode(query).tolist()
    
    prod_result = qdrant_client.query_points(
        collection_name=QDRANT_PRODUCT_COLLECTION,
        query=vector,
        limit=1 # Lấy sản phẩm sát nhất
    )
    
    if not prod_result.points:
        return "Xin lỗi, tôi không thể tìm thấy sản phẩm nào khớp với câu hỏi của bạn để lấy đánh giá."
        
    top_product = prod_result.points[0].payload
    product_id = top_product.get('id')
    product_name = top_product.get('name')
    
    # Dùng product_id đã tìm được đập vào bảng ratings để filter
    filter_condition = Filter(
        must=[FieldCondition(key="product_id", match=MatchValue(value=product_id))]
    )
    
    rating_result = qdrant_client.query_points(
        collection_name=QDRANT_RATING_COLLECTION,
        query=vector, # Ưu tiên lấy rating sát nghĩa query (ví dụ khen/chê)
        query_filter=filter_condition,
        limit=5
    )
    
    context_chunks = []
    if not rating_result.points:
        context_chunks.append("Chưa có đánh giá nào cho sản phẩm này.")
    else:
        for hit in rating_result.points:
            r = hit.payload
            chunk = f"- Khách {r.get('user', {}).get('name', 'Ẩn danh')} ({r.get('star')} sao): {r.get('content', '')}"
            context_chunks.append(chunk)
    
    context_text = f"Sản phẩm đang được hỏi là: {product_name}\n" + "\n".join(context_chunks)
    
    # Gọi LLM
    prompt_str = REVIEW_SUMMARY_PROMPT.format(context=context_text, query=query)
    response = llm.invoke(prompt_str)
    
    return response.content
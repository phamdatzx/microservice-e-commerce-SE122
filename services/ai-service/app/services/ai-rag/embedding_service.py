from functools import lru_cache
from typing import List

from langchain.text_splitter import RecursiveCharacterTextSplitter
from sentence_transformers import SentenceTransformer

from app.core.config import get_settings
from app.schemas.product_doc_chunk import ProductDocChunkRequest
from app.schemas.product_vector import ProductPayload
from app.schemas.review_doc_chunk import ReviewDocChunkRequest


@lru_cache
def get_embedding_model() -> SentenceTransformer:
    settings = get_settings()
    return SentenceTransformer(settings.EMBEDDING_MODEL_NAME)


def _build_product_text(payload: ProductPayload) -> str:
    """
    Build a textual representation of the product for embedding.
    You can tune this later (add description, attributes, etc.).
    """
    parts: List[str] = [
        payload.name,
        payload.category_name,
    ]

    return " | ".join(parts)


def compute_product_embedding(payload: ProductPayload) -> List[float]:
    model = get_embedding_model()
    text = _build_product_text(payload)
    vector = model.encode(text, normalize_embeddings=True)
    return vector.tolist()


def chunk_product_description(description: str, product_name: str, category_name: str) -> List[str]:
    """
    Chunk the product description into smaller pieces for embedding.
    Includes product name and category in the text for better context.
    """
    full_text = f"{product_name} | {category_name} | {description}"
    splitter = RecursiveCharacterTextSplitter(
        chunk_size=500,  # ~200-500 tokens depending on language
        chunk_overlap=50,
        separators=["\n\n", "\n", ". ", " ", ""]
    )
    return splitter.split_text(full_text)


def compute_product_doc_chunks(request: ProductDocChunkRequest) -> List[float]:
    """
    Compute embedding for a single product doc chunk.
    """
    model = get_embedding_model()
    vector = model.encode(request.chunk_text, normalize_embeddings=True)
    return vector.tolist()


def compute_review_doc_chunks(request: ReviewDocChunkRequest) -> List[float]:
    """
    Compute embedding for a single review doc chunk.
    """
    model = get_embedding_model()
    vector = model.encode(request.chunk_text, normalize_embeddings=True)
    return vector.tolist()


def generate_product_doc_chunks(product_id: str, name: str, description: str, category_names: List[str], seller_id: str, category_ids: List[str]) -> List[ProductDocChunkRequest]:
    """
    Generate chunk requests from full product data.
    """
    full_text = f"{name} | {' | '.join(category_names)} | {description}"
    chunks = chunk_text(full_text)
    requests = []
    for i, chunk in enumerate(chunks):
        request = ProductDocChunkRequest(
            id=f"{product_id}_chunk_{i}",
            product_id=product_id,
            chunk_text=chunk,
            chunk_index=i,
            category_ids=category_ids,
            seller_id=seller_id
        )
        requests.append(request)
    return requests


def generate_review_doc_chunks(rating_id: str, product_id: str, user_id: str, star: int, content: str, created_at: str) -> List[ReviewDocChunkRequest]:
    """
    Generate chunk requests from full review data.
    """
    full_text = f"{star} sao: {content}"
    chunks = chunk_text(full_text, chunk_size=300, overlap=30)
    requests = []
    for i, chunk in enumerate(chunks):
        request = ReviewDocChunkRequest(
            id=f"{rating_id}_chunk_{i}",
            product_id=product_id,
            rating_id=rating_id,
            user_id=user_id,
            star=star,
            chunk_text=chunk,
            chunk_index=i,
            created_at=created_at
        )
        requests.append(request)
    return requests


def chunk_text(text: str, chunk_size: int = 500, overlap: int = 50) -> List[str]:
    """
    Generic text chunking function.
    """
    splitter = RecursiveCharacterTextSplitter(
        chunk_size=chunk_size,
        chunk_overlap=overlap,
        separators=["\n\n", "\n", ". ", " ", ""]
    )
    return splitter.split_text(text)


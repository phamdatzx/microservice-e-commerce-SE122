from typing import List

from pydantic import BaseModel


class ReviewDocChunk(BaseModel):
    """Chunked document for product reviews."""

    id: str  # Unique chunk ID (e.g., rating_id + chunk_index)
    product_id: str
    rating_id: str
    user_id: str
    star: int
    chunk_text: str
    chunk_index: int  # For ordering if review is chunked
    created_at: str  # ISO string
    source_type: str = "review"  # Fixed for this schema


class ReviewDocChunkRequest(BaseModel):
    """Incoming request for upserting a review doc chunk."""

    id: str
    product_id: str
    rating_id: str
    user_id: str
    star: int
    chunk_text: str
    chunk_index: int
    created_at: str


class ReviewDocChunkVector(BaseModel):
    """Internal representation with computed vector."""

    id: str
    vector: List[float]
    payload: ReviewDocChunk</content>
<parameter name="filePath">c:\Users\nguye\Desktop\School\Đồ án 2\microservice-e-commerce-SE122\services\ai-service\app\schemas\review_doc_chunk.py
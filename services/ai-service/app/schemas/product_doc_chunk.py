from typing import List

from pydantic import BaseModel


class ProductDocChunk(BaseModel):
    """Chunked document for product details (e.g., description, specs)."""

    id: str  # Unique chunk ID (e.g., product_id + chunk_index)
    product_id: str
    chunk_text: str
    chunk_index: int  # For ordering chunks
    category_ids: List[str]
    seller_id: str
    source_type: str = "product_detail"  # Fixed for this schema


class ProductDocChunkRequest(BaseModel):
    """Incoming request for upserting a product doc chunk."""

    id: str
    product_id: str
    chunk_text: str
    chunk_index: int
    category_ids: List[str]
    seller_id: str


class ProductDocChunkVector(BaseModel):
    """Internal representation with computed vector."""

    id: str
    vector: List[float]
    payload: ProductDocChunk</content>
<parameter name="filePath">c:\Users\nguye\Desktop\School\Đồ án 2\microservice-e-commerce-SE122\services\ai-service\app\schemas\product_doc_chunk.py
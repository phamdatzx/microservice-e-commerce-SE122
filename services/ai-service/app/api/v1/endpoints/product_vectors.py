from fastapi import APIRouter, status

from app.schemas.product_vector import (
    ProductIndexRequest,
    ProductVector,
    UpsertProductVectorResponse,
)
from app.services.embedding_service import compute_product_embedding
from app.services.qdrant_service import upsert_product_vector

router = APIRouter()


@router.post(
    "/products/embedding",
    response_model=UpsertProductVectorResponse,
    status_code=status.HTTP_201_CREATED,
)
async def add_product_embedding(
    product: ProductIndexRequest,
) -> UpsertProductVectorResponse:
    """
    Insert or update a product in Qdrant.
    The embedding vector is computed using the BGE-M3 model.
    """
    vector = compute_product_embedding(product.payload)
    vector_product = ProductVector(id=product.id, vector=vector, payload=product.payload)
    upsert_product_vector(vector_product)
    return UpsertProductVectorResponse(id=product.id, status="stored")


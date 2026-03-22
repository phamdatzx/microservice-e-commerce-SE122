from typing import List

from pydantic import BaseModel


class ProductPayload(BaseModel):
    id: str
    name: str
    category_id: str
    category_name: str
    seller_id: str


class ProductIndexRequest(BaseModel):
    """Incoming request body from API (no vector field)."""

    id: str
    payload: ProductPayload


class ProductVector(BaseModel):
    """Internal representation including computed vector."""

    id: str
    vector: List[float]
    payload: ProductPayload


class UpsertProductVectorResponse(BaseModel):
    id: str
    status: str

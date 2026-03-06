from typing import List

from pydantic import BaseModel


class ProductPayload(BaseModel):
    name: str
    price_min: int
    price_max: int
    rating: float
    sold_count: int
    seller_id: str
    category_ids: List[str]
    category_names: List[str]


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

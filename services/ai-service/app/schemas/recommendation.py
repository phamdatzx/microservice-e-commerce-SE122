from typing import Any, List

from pydantic import BaseModel, Field


class RecommendationRequest(BaseModel):
    user_id: str = Field(..., description="User id to base recommendations on")
    limit: int = Field(10, gt=0, le=100, description="Max number of products to return")


class RecommendedProduct(BaseModel):
    id: str
    score: float
    payload: dict[str, Any] | None = None


class RecommendationResponse(BaseModel):
    user_id: str
    items: List[RecommendedProduct]


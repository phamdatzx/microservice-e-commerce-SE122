from typing import List

from pydantic import BaseModel, Field


class CFSimilarItem(BaseModel):
    product_id: str
    score: float


class CFRecommendationResponse(BaseModel):
    product_id: str
    similar_items: List[CFSimilarItem] = Field(default_factory=list)

from typing import List

from pydantic import BaseModel, Field


class UserProductWeight(BaseModel):
    product_id: str
    weight: float = Field(gt=0, description="Relative weight for this product")


class UserVectorRequest(BaseModel):
    user_id: str
    items: List[UserProductWeight]


class UserVectorResponse(BaseModel):
    user_id: str
    status: str

from fastapi import APIRouter, HTTPException, status

from app.schemas.recommendation import (
    RecommendationRequest,
    RecommendationResponse,
    RecommendedProduct,
)
from app.services.recommendation_service import recommend_products_for_user

router = APIRouter()


@router.post(
    "/users/recommendations",
    response_model=RecommendationResponse,
    status_code=status.HTTP_200_OK,
)
async def recommend_products(body: RecommendationRequest) -> RecommendationResponse:
    """
    Recommend products based on a stored user vector.
    """
    try:
        scored_points = recommend_products_for_user(
            user_id=body.user_id,
            limit=body.limit,
        )
    except ValueError as exc:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=str(exc),
        ) from exc

    items = [
        RecommendedProduct(
            id=str(p.id),
            score=float(p.score),
            payload=p.payload,
        )
        for p in scored_points
    ]

    # If there are simply no matches, return 200 with empty items (as requested)
    return RecommendationResponse(user_id=body.user_id, items=items)


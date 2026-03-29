from fastapi import APIRouter, HTTPException, Query, status

from app.schemas.cf_recommendation import CFRecommendationResponse, CFSimilarItem
from app.services.cf_service import get_cf_recommendations

router = APIRouter()


@router.get(
    "/recommend/cf/{product_id}",
    response_model=CFRecommendationResponse,
    status_code=status.HTTP_200_OK,
)
async def cf_recommend(
    product_id: str,
    limit: int = Query(10, gt=0, le=100, description="Max similar items to return"),
) -> CFRecommendationResponse:
    """
    Return items similar to the given product based on collaborative filtering.

    Results are pre-computed by the CF batch worker and stored in MongoDB.
    """
    raw = get_cf_recommendations(product_id, limit=limit)

    if not raw:
        raise HTTPException(
            status_code=status.HTTP_404_NOT_FOUND,
            detail=f"No CF similarities found for product_id={product_id}",
        )

    items = [CFSimilarItem(**item) for item in raw]
    return CFRecommendationResponse(product_id=product_id, similar_items=items)

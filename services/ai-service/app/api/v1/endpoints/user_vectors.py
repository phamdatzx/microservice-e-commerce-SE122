from fastapi import APIRouter, HTTPException, status

from app.schemas.user_vector import UserVectorRequest, UserVectorResponse
from app.services.user_vector_service import compute_user_vector, upsert_user_vector

router = APIRouter()


@router.post(
    "/users/vector",
    response_model=UserVectorResponse,
    status_code=status.HTTP_201_CREATED,
)
async def upsert_user_embedding(body: UserVectorRequest) -> UserVectorResponse:
    """
    Compute and store/update a user vector in Qdrant,
    based on a list of (product_id, weight).
    """
    if not body.items:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail="items must not be empty",
        )

    try:
        vector = compute_user_vector(body.items)
        upsert_user_vector(body.user_id, vector)
    except ValueError as exc:
        raise HTTPException(
            status_code=status.HTTP_400_BAD_REQUEST,
            detail=str(exc),
        ) from exc

    return UserVectorResponse(user_id=body.user_id, status="stored")


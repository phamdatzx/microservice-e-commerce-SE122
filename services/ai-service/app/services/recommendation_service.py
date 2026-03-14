from typing import List

from qdrant_client.http import models as qmodels

from app.core.config import get_settings
from app.services.qdrant_service import get_qdrant_client


def _get_user_vector(user_id: str) -> List[float]:
    settings = get_settings()
    client = get_qdrant_client()

    result = client.retrieve(
        collection_name=settings.QDRANT_USER_COLLECTION_NAME,
        ids=[user_id],
        with_vectors=True,
        with_payload=False,
    )

    if not result:
        raise ValueError(f"User vector not found for user_id={user_id}")

    point = result[0]
    if point.vector is None:
        raise ValueError(f"User vector is empty for user_id={user_id}")

    return point.vector


def recommend_products_for_user(user_id: str, limit: int) -> list[qmodels.ScoredPoint]:
    settings = get_settings()
    client = get_qdrant_client()

    user_vector = _get_user_vector(user_id)

    search_result = client.search(
        collection_name=settings.QDRANT_COLLECTION_NAME,
        query_vector=user_vector,
        limit=limit,
        with_payload=True,
    )

    return search_result


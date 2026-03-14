from typing import List

import numpy as np
from qdrant_client.http import models as qmodels

from app.core.config import get_settings
from app.schemas.user_vector import UserProductWeight
from app.services.qdrant_service import _ensure_collection, get_qdrant_client


def _get_product_vectors(product_ids: List[str]) -> List[List[float]]:
    settings = get_settings()
    client = get_qdrant_client()

    if not product_ids:
        return []

    result = client.retrieve(
        collection_name=settings.QDRANT_COLLECTION_NAME,
        ids=product_ids,
        with_vectors=True,
        with_payload=False,
    )

    # Preserve order of requested product_ids as much as possible
    id_to_vec = {point.id: point.vector for point in result}
    return [id_to_vec[p_id] for p_id in product_ids if p_id in id_to_vec]


def compute_user_vector(items: List[UserProductWeight]) -> List[float]:
    product_ids = [item.product_id for item in items]
    vectors = _get_product_vectors(product_ids)

    if not vectors:
        raise ValueError("No product vectors found for given product_ids")

    if len(vectors) != len(items):
        raise ValueError("Some product_ids do not exist in the products collection")

    weights = np.array([item.weight for item in items], dtype=float)
    vectors_np = np.array(vectors, dtype=float)

    if weights.sum() <= 0:
        raise ValueError("Sum of weights must be positive")

    # Weighted average of product vectors
    weighted = (vectors_np.T * weights).T
    user_vec = weighted.sum(axis=0) / weights.sum()

    # Normalize for cosine similarity
    norm = np.linalg.norm(user_vec)
    if norm == 0:
        raise ValueError("Resulting user vector has zero norm")

    return (user_vec / norm).tolist()


def upsert_user_vector(user_id: str, vector: List[float]) -> None:
    settings = get_settings()
    client = get_qdrant_client()

    dim = len(vector)
    if dim == 0:
        raise ValueError("User vector must not be empty")

    _ensure_collection(client, settings.QDRANT_USER_COLLECTION_NAME, dim)

    client.upsert(
        collection_name=settings.QDRANT_USER_COLLECTION_NAME,
        points=[
            qmodels.PointStruct(
                id=user_id,
                vector=vector,
                payload={"user_id": user_id},
            )
        ],
    )


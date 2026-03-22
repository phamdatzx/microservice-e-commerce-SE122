import logging
from collections import defaultdict
from typing import Any, List

import numpy as np
from qdrant_client.http import models as qmodels

from app.core.config import get_settings
from app.schemas.user_vector import UserProductWeight
from app.services.qdrant_service import _ensure_collection, get_qdrant_client

logger = logging.getLogger(__name__)


def _get_product_vector_map(product_ids: List[str]) -> dict[str, List[float]]:
    settings = get_settings()
    client = get_qdrant_client()

    if not product_ids:
        return {}

    result = client.retrieve(
        collection_name=settings.QDRANT_COLLECTION_NAME,
        ids=product_ids,
        with_vectors=True,
        with_payload=False,
    )

    out: dict[str, List[float]] = {}
    for point in result:
        if point.vector is not None:
            out[point.id] = point.vector
    return out


def _weighted_average_product_vectors(
    items: List[UserProductWeight], vec_map: dict[str, List[float]]
) -> List[float]:
    weights: list[float] = []
    vectors: list[List[float]] = []
    for item in items:
        vec = vec_map.get(item.product_id)
        if vec is None:
            raise ValueError(f"Product vector missing for product_id={item.product_id!r}")
        weights.append(float(item.weight))
        vectors.append(vec)

    weights_np = np.array(weights, dtype=float)
    vectors_np = np.array(vectors, dtype=float)

    if weights_np.sum() <= 0:
        raise ValueError("Sum of weights must be positive")

    weighted = (vectors_np.T * weights_np).T
    user_vec = weighted.sum(axis=0) / weights_np.sum()

    norm = np.linalg.norm(user_vec)
    if norm == 0:
        raise ValueError("Resulting user vector has zero norm")

    return (user_vec / norm).tolist()


def compute_user_vector(items: List[UserProductWeight]) -> List[float]:
    product_ids = [item.product_id for item in items]
    unique_ids = list(dict.fromkeys(product_ids))
    vec_map = _get_product_vector_map(unique_ids)

    if not vec_map:
        raise ValueError("No product vectors found for given product_ids")

    missing = [pid for pid in unique_ids if pid not in vec_map]
    if missing:
        raise ValueError(
            "Some product_ids do not exist in the products collection: "
            + ", ".join(missing[:10])
            + ("..." if len(missing) > 10 else "")
        )

    return _weighted_average_product_vectors(items, vec_map)


def compute_user_vector_from_interaction_docs(
    interactions: list[dict[str, Any]],
) -> List[float]:
    """
    Build a user vector from Mongo interaction documents (newest-first order).
    Aggregates scores per product_id, then weighted-averages product vectors.
    Skips product_ids missing from Qdrant (with a warning).
    """
    agg: defaultdict[str, float] = defaultdict(float)
    for doc in interactions:
        pid = str(doc.get("product_id", ""))
        if not pid:
            continue
        agg[pid] += float(doc.get("score", 0.0))

    items = [
        UserProductWeight(product_id=pid, weight=w)
        for pid, w in agg.items()
        if w > 0
    ]
    if not items:
        raise ValueError("No positive weights after aggregating interactions")

    vec_map = _get_product_vector_map([i.product_id for i in items])
    kept_items: list[UserProductWeight] = []
    for it in items:
        if it.product_id in vec_map:
            kept_items.append(it)
        else:
            logger.warning(
                "Skipping product_id %s (no vector in Qdrant products collection)",
                it.product_id,
            )

    if not kept_items:
        raise ValueError(
            "No product vectors found for aggregated interactions (all products missing in Qdrant)"
        )

    return _weighted_average_product_vectors(kept_items, vec_map)


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


"""
Item-based Collaborative Filtering service.

Pipeline:
  1. Load recent interactions from MongoDB
  2. Build sparse user-item matrix
  3. Compute item-item cosine similarity row-by-row (memory safe)
  4. Store top-K similar items per product in MongoDB
  5. Query similar items for a given product
"""

from __future__ import annotations

import logging
from datetime import datetime, timedelta, timezone
from typing import Any

import numpy as np
import pandas as pd
from pymongo import ASCENDING, UpdateOne
from scipy.sparse import csr_matrix
from sklearn.metrics.pairwise import cosine_similarity

from app.core.config import get_settings
from app.services.interaction_service import get_interactions_collection, get_mongo_client

logger = logging.getLogger(__name__)


# ---------------------------------------------------------------------------
# MongoDB helpers
# ---------------------------------------------------------------------------


def _get_similarity_collection():
    settings = get_settings()
    db = get_mongo_client()[settings.MONGO_DB_NAME]
    return db[settings.ITEM_SIMILARITY_COLLECTION_NAME]


# ---------------------------------------------------------------------------
# 1. Load interactions
# ---------------------------------------------------------------------------


def load_interactions() -> pd.DataFrame:
    """
    Load interaction data from MongoDB for the last 30 days.

    Applies filters:
      - Only last 30 days
      - Drop users with < CF_MIN_USER_INTERACTIONS
      - Drop items with < CF_MIN_ITEM_INTERACTIONS
      - Keep only top CF_MAX_PRODUCTS by interaction count
    """
    settings = get_settings()
    col = get_interactions_collection()

    cutoff = datetime.now(timezone.utc) - timedelta(days=30)

    cursor = col.find(
        {"timestamp": {"$gte": cutoff}},
        {"_id": 0, "user_id": 1, "product_id": 1, "score": 1},
    )

    records = list(cursor)
    if not records:
        logger.warning("No interactions found in the last 30 days")
        return pd.DataFrame(columns=["user_id", "product_id", "score"])

    df = pd.DataFrame(records)
    logger.info("Loaded %d raw interactions from last 30 days", len(df))

    # --- Filter users with too few interactions ---
    user_counts = df["user_id"].value_counts()
    valid_users = user_counts[user_counts >= settings.CF_MIN_USER_INTERACTIONS].index
    df = df[df["user_id"].isin(valid_users)]
    logger.info(
        "After user filter (>=%d interactions): %d interactions",
        settings.CF_MIN_USER_INTERACTIONS,
        len(df),
    )

    # --- Filter items with too few interactions ---
    item_counts = df["product_id"].value_counts()
    valid_items = item_counts[item_counts >= settings.CF_MIN_ITEM_INTERACTIONS].index
    df = df[df["product_id"].isin(valid_items)]
    logger.info(
        "After item filter (>=%d interactions): %d interactions",
        settings.CF_MIN_ITEM_INTERACTIONS,
        len(df),
    )

    # --- Keep only top-N products by interaction count ---
    item_counts = df["product_id"].value_counts()
    top_items = item_counts.head(settings.CF_MAX_PRODUCTS).index
    df = df[df["product_id"].isin(top_items)]
    logger.info(
        "After top-%d products cap: %d interactions, %d unique products",
        settings.CF_MAX_PRODUCTS,
        len(df),
        df["product_id"].nunique(),
    )

    return df.reset_index(drop=True)


# ---------------------------------------------------------------------------
# 2. Build sparse user-item matrix
# ---------------------------------------------------------------------------


def build_user_item_matrix(
    df: pd.DataFrame,
) -> tuple[csr_matrix, np.ndarray, np.ndarray]:
    """
    Build a sparse user × item matrix from interaction data.

    Returns:
        (sparse_matrix, user_ids_array, product_ids_array)
    """
    # Aggregate scores: if a user interacted with the same product multiple
    # times, sum the scores.
    agg = df.groupby(["user_id", "product_id"], as_index=False)["score"].sum()

    user_ids = agg["user_id"].astype("category")
    product_ids = agg["product_id"].astype("category")

    row = user_ids.cat.codes.values
    col = product_ids.cat.codes.values
    data = agg["score"].values.astype(np.float32)

    sparse_mat = csr_matrix(
        (data, (row, col)),
        shape=(user_ids.cat.categories.size, product_ids.cat.categories.size),
    )

    logger.info(
        "Built sparse matrix: %d users × %d items, nnz=%d",
        sparse_mat.shape[0],
        sparse_mat.shape[1],
        sparse_mat.nnz,
    )

    return sparse_mat, np.array(user_ids.cat.categories), np.array(product_ids.cat.categories)


# ---------------------------------------------------------------------------
# 3. Compute item similarity (row-by-row, memory safe)
# ---------------------------------------------------------------------------


def compute_item_similarity(
    sparse_mat: csr_matrix,
    product_ids: np.ndarray,
    top_k: int | None = None,
) -> dict[str, list[dict[str, Any]]]:
    """
    Compute cosine similarity per item against all other items.
    Never materializes the full N×N matrix.

    Returns:
        { product_id: [{ "product_id": ..., "score": ... }, ...] }
    """
    settings = get_settings()
    if top_k is None:
        top_k = settings.CF_TOP_K

    # Transpose to item × user
    item_user_matrix = sparse_mat.T.tocsr()
    n_items = item_user_matrix.shape[0]

    results: dict[str, list[dict[str, Any]]] = {}

    for i in range(n_items):
        # Compute similarity of item i against ALL items → 1×N vector
        item_vec = item_user_matrix[i]
        sim_row = cosine_similarity(item_vec, item_user_matrix).flatten()

        # Exclude self
        sim_row[i] = -1.0

        # Get top-K indices
        if top_k >= n_items:
            top_indices = np.argsort(sim_row)[::-1][: n_items - 1]
        else:
            # Use argpartition for efficiency on large arrays
            top_indices = np.argpartition(sim_row, -top_k)[-top_k:]
            top_indices = top_indices[np.argsort(sim_row[top_indices])[::-1]]

        # Filter out non-positive similarities
        similar_items = []
        for idx in top_indices:
            score = float(sim_row[idx])
            if score <= 0:
                break
            similar_items.append(
                {"product_id": str(product_ids[idx]), "score": round(score, 6)}
            )

        pid = str(product_ids[i])
        if similar_items:
            results[pid] = similar_items

        if (i + 1) % 500 == 0:
            logger.info("Computed similarity for %d / %d items", i + 1, n_items)

    logger.info(
        "Computed similarities for %d items (with at least 1 similar item)", len(results)
    )

    return results


# ---------------------------------------------------------------------------
# 4. Store results in MongoDB
# ---------------------------------------------------------------------------


def store_item_similarities(
    similarities: dict[str, list[dict[str, Any]]],
) -> None:
    """
    Upsert item similarity results into MongoDB.
    Removes stale entries not in the current batch.
    """
    col = _get_similarity_collection()

    # Ensure index
    col.create_index([("product_id", ASCENDING)], unique=True)

    now = datetime.now(timezone.utc)
    operations = []
    for product_id, similar_items in similarities.items():
        operations.append(
            UpdateOne(
                {"product_id": product_id},
                {
                    "$set": {
                        "product_id": product_id,
                        "similar_items": similar_items,
                        "updated_at": now,
                    }
                },
                upsert=True,
            )
        )

    if operations:
        # Write in batches of 1000
        batch_size = 1000
        for start in range(0, len(operations), batch_size):
            batch = operations[start : start + batch_size]
            result = col.bulk_write(batch, ordered=False)
            logger.info(
                "Bulk write batch %d-%d: upserted=%d, modified=%d",
                start,
                start + len(batch),
                result.upserted_count,
                result.modified_count,
            )

    # Remove stale products not in this batch
    current_product_ids = list(similarities.keys())
    delete_result = col.delete_many({"product_id": {"$nin": current_product_ids}})
    if delete_result.deleted_count:
        logger.info("Removed %d stale similarity entries", delete_result.deleted_count)


# ---------------------------------------------------------------------------
# 5. Query similar items (used by API)
# ---------------------------------------------------------------------------


def get_cf_recommendations(
    product_id: str, limit: int = 10
) -> list[dict[str, Any]]:
    """
    Return pre-computed similar items for a product from MongoDB.
    """
    col = _get_similarity_collection()
    doc = col.find_one({"product_id": product_id}, {"_id": 0})

    if not doc or not doc.get("similar_items"):
        return []

    items = doc["similar_items"]
    return items[:limit]

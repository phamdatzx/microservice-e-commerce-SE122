"""
Batch worker: compute item-based collaborative filtering similarities.

Usage:
    python -m app.cf_worker

This is a one-shot job (not a long-running consumer).
Schedule via cron or run manually.
"""

from __future__ import annotations

import logging
import time

from app.core.config import get_settings
from app.services.cf_service import (
    build_user_item_matrix,
    compute_item_similarity,
    load_interactions,
    store_item_similarities,
)

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO, format="%(asctime)s %(levelname)s %(name)s: %(message)s")


def main() -> None:
    settings = get_settings()

    if not settings.MONGO_URI or not settings.MONGO_DB_NAME:
        raise RuntimeError(
            "Set MONGO_URI and MONGO_DB_NAME for the CF worker."
        )

    logger.info("=== CF batch job started ===")
    t0 = time.perf_counter()

    # Step 1: Load interactions
    logger.info("Step 1/4: Loading interactions...")
    df = load_interactions()
    if df.empty:
        logger.warning("No interactions to process. Exiting.")
        return

    # Step 2: Build sparse matrix
    logger.info("Step 2/4: Building user-item matrix...")
    sparse_mat, user_ids, product_ids = build_user_item_matrix(df)

    # Step 3: Compute similarities
    logger.info("Step 3/4: Computing item-item similarities (top_k=%d)...", settings.CF_TOP_K)
    similarities = compute_item_similarity(sparse_mat, product_ids, top_k=settings.CF_TOP_K)

    # Step 4: Store results
    logger.info("Step 4/4: Storing similarities to MongoDB...")
    store_item_similarities(similarities)

    elapsed = time.perf_counter() - t0
    logger.info(
        "=== CF batch job completed in %.2f seconds. %d products with similarities. ===",
        elapsed,
        len(similarities),
    )


if __name__ == "__main__":
    main()

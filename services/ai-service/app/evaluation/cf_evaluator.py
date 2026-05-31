"""
Offline evaluation of the item-based CF system using MovieLens 100k.

Mapping:
    movie_id  → product_id
    user_id   → user_id
    rating    → score  (1–5 explicit ratings; items rated >= RELEVANCE_THRESHOLD are "relevant")

Pipeline per fold:
    1. Load base (train) split  → DataFrame(user_id, product_id, score)
    2. Build sparse user×item matrix  (reuses cf_service logic)
    3. Compute item-item cosine similarity  (reuses cf_service logic)
    4. For every test user: expand their training items through similarity, rank candidates
    5. Compute Precision@K, Recall@K, NDCG@K, Hit-Rate@K
    6. Average across all folds

Usage (standalone, no MongoDB required):
    python -m app.evaluation.cf_evaluator
    python -m app.evaluation.cf_evaluator --data-dir ml-100k --top-k 20 --k-values 5,10,20
"""

from __future__ import annotations

import argparse
import logging
import math
import sys
import time
from collections import defaultdict
from pathlib import Path
from typing import NamedTuple

import numpy as np
import pandas as pd
from scipy.sparse import csr_matrix
from sklearn.metrics.pairwise import cosine_similarity

logger = logging.getLogger(__name__)

# ---------------------------------------------------------------------------
# Constants / defaults
# ---------------------------------------------------------------------------

RELEVANCE_THRESHOLD = 4       # ratings >= this are considered "relevant"
DEFAULT_TOP_K = 20            # number of similar items stored per item
DEFAULT_K_VALUES = (5, 10, 20)
FOLD_PAIRS = [
    ("u1.base", "u1.test"),
    ("u2.base", "u2.test"),
    ("u3.base", "u3.test"),
    ("u4.base", "u4.test"),
    ("u5.base", "u5.test"),
]


# ---------------------------------------------------------------------------
# 1. Data loading
# ---------------------------------------------------------------------------


def load_ml100k_split(base_path: Path, test_path: Path) -> tuple[pd.DataFrame, pd.DataFrame]:
    """
    Load a single MovieLens 100k train/test split.

    Columns returned: user_id (str), product_id (str), score (float)
    """
    cols = ["user_id", "product_id", "score", "timestamp"]

    def _read(p: Path) -> pd.DataFrame:
        df = pd.read_csv(p, sep="\t", names=cols, usecols=["user_id", "product_id", "score"])
        df["user_id"] = df["user_id"].astype(str)
        df["product_id"] = df["product_id"].astype(str)
        df["score"] = df["score"].astype(float)
        return df

    return _read(base_path), _read(test_path)


# ---------------------------------------------------------------------------
# 2. Build sparse matrix (inline — no MongoDB dependency)
# ---------------------------------------------------------------------------


def build_user_item_matrix(
    df: pd.DataFrame,
) -> tuple[csr_matrix, np.ndarray, np.ndarray]:
    """
    Build a sparse user × item matrix from interactions.
    Duplicate (user, product) pairs are summed.

    Returns:
        (sparse_matrix, user_ids_array, product_ids_array)
    """
    agg = df.groupby(["user_id", "product_id"], as_index=False)["score"].sum()

    user_cat = agg["user_id"].astype("category")
    item_cat = agg["product_id"].astype("category")

    row = user_cat.cat.codes.values
    col = item_cat.cat.codes.values
    data = agg["score"].values.astype(np.float32)

    sparse_mat = csr_matrix(
        (data, (row, col)),
        shape=(user_cat.cat.categories.size, item_cat.cat.categories.size),
    )

    return sparse_mat, np.array(user_cat.cat.categories), np.array(item_cat.cat.categories)


# ---------------------------------------------------------------------------
# 3. Compute item-item cosine similarity
# ---------------------------------------------------------------------------


def compute_item_similarity(
    sparse_mat: csr_matrix,
    product_ids: np.ndarray,
    top_k: int,
) -> dict[str, list[tuple[str, float]]]:
    """
    Compute top-K cosine-similar items for every item in the matrix.
    Memory-safe: processes one row at a time.

    Returns:
        { product_id: [(similar_product_id, score), ...] }   sorted desc by score
    """
    item_user_mat = sparse_mat.T.tocsr()
    n_items = item_user_mat.shape[0]
    results: dict[str, list[tuple[str, float]]] = {}

    for i in range(n_items):
        item_vec = item_user_mat[i]
        sim_row = cosine_similarity(item_vec, item_user_mat).flatten()
        sim_row[i] = -1.0  # exclude self

        k = min(top_k, n_items - 1)
        top_indices = np.argpartition(sim_row, -k)[-k:]
        top_indices = top_indices[np.argsort(sim_row[top_indices])[::-1]]

        similar: list[tuple[str, float]] = []
        for idx in top_indices:
            score = float(sim_row[idx])
            if score <= 0:
                break
            similar.append((str(product_ids[idx]), score))

        pid = str(product_ids[i])
        if similar:
            results[pid] = similar

        if (i + 1) % 200 == 0:
            logger.debug("Similarity: %d / %d items done", i + 1, n_items)

    return results


# ---------------------------------------------------------------------------
# 4. Generate top-N recommendations for a user
# ---------------------------------------------------------------------------


def generate_user_recommendations(
    user_training_items: list[tuple[str, float]],
    item_similarities: dict[str, list[tuple[str, float]]],
    top_n: int,
    exclude_seen: bool = True,
) -> list[str]:
    """
    Given a user's training interactions and the precomputed item similarity map,
    produce a ranked list of top-N recommended item IDs.

    Strategy: for each training item, weight its similar items by
        candidate_score += training_score * similarity_score
    Then rank candidates descending, exclude already-seen items.

    Args:
        user_training_items: [(product_id, score), ...]
        item_similarities:   { product_id: [(similar_id, sim_score), ...] }
        top_n:               number of items to return
        exclude_seen:        whether to exclude items the user already interacted with

    Returns:
        list of product_id strings, length <= top_n
    """
    seen = {pid for pid, _ in user_training_items} if exclude_seen else set()
    candidate_scores: dict[str, float] = defaultdict(float)

    for pid, train_score in user_training_items:
        neighbors = item_similarities.get(pid, [])
        for neighbor_id, sim_score in neighbors:
            if neighbor_id not in seen:
                candidate_scores[neighbor_id] += train_score * sim_score

    if not candidate_scores:
        return []

    ranked = sorted(candidate_scores.items(), key=lambda x: x[1], reverse=True)
    return [pid for pid, _ in ranked[:top_n]]


# ---------------------------------------------------------------------------
# 5. Evaluation metrics
# ---------------------------------------------------------------------------


class MetricResult(NamedTuple):
    precision: float
    recall: float
    ndcg: float
    hit_rate: float
    n_users: int


def _dcg(relevances: list[int]) -> float:
    return sum(rel / math.log2(rank + 2) for rank, rel in enumerate(relevances))


def compute_metrics_at_k(
    recommendations: dict[str, list[str]],
    ground_truth: dict[str, set[str]],
    k: int,
) -> MetricResult:
    """
    Compute Precision@K, Recall@K, NDCG@K, Hit-Rate@K.

    Args:
        recommendations: { user_id: [ranked product_ids, top-N] }
        ground_truth:    { user_id: {relevant product_ids} }
        k:               cutoff rank

    Returns:
        MetricResult with averaged metrics over users that have ground-truth
    """
    precisions, recalls, ndcgs, hits = [], [], [], []

    for user_id, relevant_items in ground_truth.items():
        if not relevant_items:
            continue

        recs = recommendations.get(user_id, [])[:k]
        if not recs:
            precisions.append(0.0)
            recalls.append(0.0)
            ndcgs.append(0.0)
            hits.append(0.0)
            continue

        hits_list = [1 if r in relevant_items else 0 for r in recs]
        n_hits = sum(hits_list)

        precisions.append(n_hits / k)
        recalls.append(n_hits / len(relevant_items))
        hits.append(1.0 if n_hits > 0 else 0.0)

        # NDCG
        ideal = [1] * min(len(relevant_items), k)
        ndcgs.append(_dcg(hits_list) / _dcg(ideal) if ideal else 0.0)

    if not precisions:
        return MetricResult(0.0, 0.0, 0.0, 0.0, 0)

    n = len(precisions)
    return MetricResult(
        precision=sum(precisions) / n,
        recall=sum(recalls) / n,
        ndcg=sum(ndcgs) / n,
        hit_rate=sum(hits) / n,
        n_users=n,
    )


# ---------------------------------------------------------------------------
# 6. Single-fold evaluation
# ---------------------------------------------------------------------------


def evaluate_fold(
    base_df: pd.DataFrame,
    test_df: pd.DataFrame,
    top_k: int,
    k_values: tuple[int, ...],
    relevance_threshold: float = RELEVANCE_THRESHOLD,
) -> dict[int, MetricResult]:
    """
    Run the full CF pipeline on one train/test split and return metrics.

    Returns:
        { k: MetricResult }
    """
    # Build matrix and compute similarities from training data
    sparse_mat, _user_ids, product_ids = build_user_item_matrix(base_df)
    item_similarities = compute_item_similarity(sparse_mat, product_ids, top_k=top_k)

    # Build per-user training items: { user_id: [(product_id, score), ...] }
    user_train: dict[str, list[tuple[str, float]]] = defaultdict(list)
    for _, row in base_df.iterrows():
        user_train[row["user_id"]].append((row["product_id"], row["score"]))

    # Build ground truth from test split: items rated >= threshold
    ground_truth: dict[str, set[str]] = defaultdict(set)
    for _, row in test_df.iterrows():
        if row["score"] >= relevance_threshold:
            ground_truth[row["user_id"]].add(row["product_id"])

    # Only evaluate users that exist in both train and test, and have relevant items
    eval_users = set(user_train.keys()) & set(ground_truth.keys())
    logger.info(
        "  Evaluating %d users (have both training items and relevant test items)",
        len(eval_users),
    )

    # Generate recommendations for each eval user (at max(k_values) depth)
    max_k = max(k_values)
    recommendations: dict[str, list[str]] = {}
    for user_id in eval_users:
        recs = generate_user_recommendations(
            user_training_items=user_train[user_id],
            item_similarities=item_similarities,
            top_n=max_k,
        )
        recommendations[user_id] = recs

    # Filter ground_truth to eval users only
    eval_ground_truth = {uid: ground_truth[uid] for uid in eval_users}

    return {k: compute_metrics_at_k(recommendations, eval_ground_truth, k) for k in k_values}


# ---------------------------------------------------------------------------
# 7. 5-fold cross-validation runner
# ---------------------------------------------------------------------------


def run_5fold_cv(
    data_dir: Path,
    top_k: int = DEFAULT_TOP_K,
    k_values: tuple[int, ...] = DEFAULT_K_VALUES,
    relevance_threshold: float = RELEVANCE_THRESHOLD,
) -> dict[int, MetricResult]:
    """
    Run 5-fold cross-validation using the standard MovieLens u1–u5 splits.

    Returns:
        { k: averaged MetricResult across all 5 folds }
    """
    fold_results: dict[int, list[MetricResult]] = defaultdict(list)

    for fold_idx, (base_name, test_name) in enumerate(FOLD_PAIRS, start=1):
        base_path = data_dir / base_name
        test_path = data_dir / test_name

        if not base_path.exists() or not test_path.exists():
            logger.warning("Fold %d: files not found (%s / %s), skipping", fold_idx, base_path, test_path)
            continue

        logger.info("=== Fold %d / %d ===", fold_idx, len(FOLD_PAIRS))
        t0 = time.perf_counter()

        base_df, test_df = load_ml100k_split(base_path, test_path)
        logger.info(
            "  Loaded: %d train interactions, %d test interactions",
            len(base_df), len(test_df),
        )

        fold_metrics = evaluate_fold(base_df, test_df, top_k, k_values, relevance_threshold)

        for k, result in fold_metrics.items():
            fold_results[k].append(result)
            logger.info(
                "  @K=%2d  P=%.4f  R=%.4f  NDCG=%.4f  HR=%.4f  (users=%d)",
                k,
                result.precision,
                result.recall,
                result.ndcg,
                result.hit_rate,
                result.n_users,
            )

        elapsed = time.perf_counter() - t0
        logger.info("  Fold %d done in %.1f s", fold_idx, elapsed)

    # Average across folds
    averaged: dict[int, MetricResult] = {}
    for k, results in fold_results.items():
        n_folds = len(results)
        averaged[k] = MetricResult(
            precision=sum(r.precision for r in results) / n_folds,
            recall=sum(r.recall for r in results) / n_folds,
            ndcg=sum(r.ndcg for r in results) / n_folds,
            hit_rate=sum(r.hit_rate for r in results) / n_folds,
            n_users=int(sum(r.n_users for r in results) / n_folds),
        )

    return averaged


# ---------------------------------------------------------------------------
# 8. Pretty-print results
# ---------------------------------------------------------------------------


def print_results(averaged: dict[int, MetricResult]) -> None:
    header = f"{'K':>4}  {'Precision@K':>12}  {'Recall@K':>10}  {'NDCG@K':>8}  {'HitRate@K':>10}  {'Avg Users':>10}"
    print("\n" + "=" * len(header))
    print("5-Fold Cross-Validation Results (MovieLens 100k — item-based CF)")
    print("=" * len(header))
    print(header)
    print("-" * len(header))
    for k in sorted(averaged.keys()):
        r = averaged[k]
        print(
            f"{k:>4}  {r.precision:>12.4f}  {r.recall:>10.4f}  {r.ndcg:>8.4f}  {r.hit_rate:>10.4f}  {r.n_users:>10}"
        )
    print("=" * len(header))


# ---------------------------------------------------------------------------
# CLI entry point
# ---------------------------------------------------------------------------


def _parse_args(argv: list[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Evaluate item-based CF on MovieLens 100k (5-fold CV)",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument(
        "--data-dir",
        type=Path,
        default=Path("ml-100k"),
        help="Path to the ml-100k directory",
    )
    parser.add_argument(
        "--top-k",
        type=int,
        default=DEFAULT_TOP_K,
        help="Number of similar items stored per item (CF_TOP_K)",
    )
    parser.add_argument(
        "--k-values",
        type=lambda s: tuple(int(x) for x in s.split(",")),
        default=",".join(str(k) for k in DEFAULT_K_VALUES),
        help="Comma-separated cutoff values, e.g. 5,10,20",
    )
    parser.add_argument(
        "--relevance-threshold",
        type=float,
        default=RELEVANCE_THRESHOLD,
        help="Minimum rating to count as relevant (e.g. 4 means ratings 4–5 are positive)",
    )
    parser.add_argument(
        "--log-level",
        choices=["DEBUG", "INFO", "WARNING"],
        default="INFO",
    )
    return parser.parse_args(argv)


def main(argv: list[str] | None = None) -> None:
    args = _parse_args(argv)
    logging.basicConfig(
        level=getattr(logging, args.log_level),
        format="%(asctime)s %(levelname)s %(name)s: %(message)s",
        stream=sys.stdout,
    )

    logger.info(
        "Starting evaluation: data_dir=%s, top_k=%d, k_values=%s, relevance_threshold=%g",
        args.data_dir,
        args.top_k,
        args.k_values,
        args.relevance_threshold,
    )

    averaged = run_5fold_cv(
        data_dir=args.data_dir,
        top_k=args.top_k,
        k_values=args.k_values,
        relevance_threshold=args.relevance_threshold,
    )

    print_results(averaged)


if __name__ == "__main__":
    main()

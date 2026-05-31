"""
Offline evaluation of the content-based filtering system.

Dataset : data-contentbased/ecommerce_products_5k.csv
          5,000 products · 20 sub-categories · 250 products each

Embed text format (mirrors the real system's embedding_service.py):
    "{name} | {sub_category}"

Ground truth:
    Products that share the same sub_category are considered relevant.
    249 relevant items exist per query product.

Metrics computed per product (query), then averaged:
    Precision@K       – fraction of top-K that share the query's sub_category
    Recall@K          – fraction of all 249 relevant items that appear in top-K
    HitRate@K         – 1 if ≥1 relevant item in top-K, else 0
    Mean Cosine Sim   – average similarity score of the top-K results

Embeddings are cached to disk (.npy + .json) so the OpenAI API is only
called once.

Usage:
    python -m app.evaluation.cb_evaluator
    python -m app.evaluation.cb_evaluator --k-values 5,10,20 --no-cache
    python -m app.evaluation.cb_evaluator --data-dir data-contentbased --cache-dir .cache/cb_eval
"""

from __future__ import annotations

import argparse
import json
import logging
import math
import sys
import time
from pathlib import Path
from typing import NamedTuple

import numpy as np
import pandas as pd
from dotenv import load_dotenv
from openai import OpenAI
from sklearn.metrics.pairwise import cosine_similarity

logger = logging.getLogger(__name__)

# ---------------------------------------------------------------------------
# Defaults
# ---------------------------------------------------------------------------

DEFAULT_DATA_DIR = Path("data-contentbased")
DEFAULT_CACHE_DIR = Path(".cache/cb_eval")
DEFAULT_CSV_NAME = "ecommerce_products_5k.csv"
DEFAULT_K_VALUES = (5, 10, 20)
EMBED_MODEL = "text-embedding-3-small"
EMBED_BATCH_SIZE = 100          # OpenAI allows up to 2048 inputs per request
RELEVANCE_COL = "sub_category"  # column that defines relevant groups


# ---------------------------------------------------------------------------
# 1. Load dataset
# ---------------------------------------------------------------------------


def load_products(data_dir: Path, csv_name: str = DEFAULT_CSV_NAME) -> pd.DataFrame:
    path = data_dir / csv_name
    if not path.exists():
        raise FileNotFoundError(f"Dataset not found: {path}")

    df = pd.read_csv(path, dtype=str)
    required = {"product_id", "name", RELEVANCE_COL}
    missing = required - set(df.columns)
    if missing:
        raise ValueError(f"CSV missing required columns: {missing}")

    df = df.reset_index(drop=True)
    logger.info(
        "Loaded %d products  |  %d sub-categories",
        len(df),
        df[RELEVANCE_COL].nunique(),
    )
    return df


def build_embed_texts(df: pd.DataFrame) -> list[str]:
    """
    Build the embedding input string for each product.
    Mirrors the real system: "{name} | {category_name}"
    """
    return (df["name"].fillna("") + " | " + df[RELEVANCE_COL].fillna("")).tolist()


# ---------------------------------------------------------------------------
# 2. Embed with OpenAI (batched + disk cache)
# ---------------------------------------------------------------------------


def _load_cache(cache_dir: Path) -> tuple[np.ndarray, list[str]] | None:
    emb_path = cache_dir / "embeddings.npy"
    ids_path = cache_dir / "product_ids.json"
    if emb_path.exists() and ids_path.exists():
        embeddings = np.load(emb_path)
        with open(ids_path) as f:
            product_ids = json.load(f)
        logger.info("Loaded cached embeddings: %s  (%d vectors)", emb_path, len(product_ids))
        return embeddings, product_ids
    return None


def _save_cache(cache_dir: Path, embeddings: np.ndarray, product_ids: list[str]) -> None:
    cache_dir.mkdir(parents=True, exist_ok=True)
    np.save(cache_dir / "embeddings.npy", embeddings)
    with open(cache_dir / "product_ids.json", "w") as f:
        json.dump(product_ids, f)
    logger.info("Embeddings cached to %s", cache_dir)


def embed_products(
    df: pd.DataFrame,
    api_key: str,
    cache_dir: Path = DEFAULT_CACHE_DIR,
    use_cache: bool = True,
    model: str = EMBED_MODEL,
) -> tuple[np.ndarray, list[str]]:
    """
    Embed all products via OpenAI text-embedding-3-small.

    Returns:
        embeddings  : np.ndarray of shape (N, 1536), L2-normalised
        product_ids : list of product_id strings aligned with rows
    """
    if use_cache:
        cached = _load_cache(cache_dir)
        if cached is not None:
            return cached

    texts = build_embed_texts(df)
    product_ids = df["product_id"].tolist()

    client = OpenAI(api_key=api_key)
    all_vectors: list[list[float]] = []

    logger.info(
        "Embedding %d products with model '%s' (batch_size=%d)...",
        len(texts), model, EMBED_BATCH_SIZE,
    )
    t0 = time.perf_counter()

    for start in range(0, len(texts), EMBED_BATCH_SIZE):
        batch = texts[start: start + EMBED_BATCH_SIZE]
        response = client.embeddings.create(model=model, input=batch)
        batch_vecs = [item.embedding for item in sorted(response.data, key=lambda x: x.index)]
        all_vectors.extend(batch_vecs)

        done = min(start + EMBED_BATCH_SIZE, len(texts))
        logger.info("  Embedded %d / %d products  (%.1fs)", done, len(texts), time.perf_counter() - t0)

    embeddings = np.array(all_vectors, dtype=np.float32)

    # L2-normalise so dot product == cosine similarity (faster later)
    norms = np.linalg.norm(embeddings, axis=1, keepdims=True)
    norms = np.where(norms == 0, 1.0, norms)
    embeddings = embeddings / norms

    logger.info(
        "Embeddings ready: shape=%s  elapsed=%.1fs",
        embeddings.shape,
        time.perf_counter() - t0,
    )

    if use_cache:
        _save_cache(cache_dir, embeddings, product_ids)

    return embeddings, product_ids


# ---------------------------------------------------------------------------
# 3. Evaluation metrics
# ---------------------------------------------------------------------------


class MetricResult(NamedTuple):
    precision: float
    recall: float
    hit_rate: float
    mean_cosine_sim: float
    n_products: int


def evaluate(
    embeddings: np.ndarray,
    product_ids: list[str],
    categories: list[str],
    k: int,
) -> MetricResult:
    """
    Evaluate content-based retrieval at cutoff K.

    For every product:
      - Rank all others by cosine similarity (dot product on normalised vectors)
      - Relevant = same sub_category (excluding self)
      - Compute Precision@K, Recall@K, HitRate@K, Mean Cosine Sim
    """
    n = len(product_ids)
    cat_arr = np.array(categories)

    precisions, recalls, hits, mean_sims = [], [], [], []

    # Process in blocks to avoid holding the full N×N matrix in memory
    block = 500
    for start in range(0, n, block):
        end = min(start + block, n)
        # (block_size, N) similarity matrix — fast because vectors are L2-normalised
        sim_block = embeddings[start:end] @ embeddings.T  # dot == cosine

        for local_i, global_i in enumerate(range(start, end)):
            sim_row = sim_block[local_i].copy()
            sim_row[global_i] = -2.0  # exclude self

            # Top-K indices by similarity
            top_k_idx = np.argpartition(sim_row, -k)[-k:]
            top_k_idx = top_k_idx[np.argsort(sim_row[top_k_idx])[::-1]]

            query_cat = categories[global_i]
            n_relevant_total = int(np.sum(cat_arr == query_cat)) - 1  # exclude self

            retrieved_cats = cat_arr[top_k_idx]
            relevance = (retrieved_cats == query_cat).astype(int)

            n_hits = relevance.sum()
            precisions.append(n_hits / k)
            recalls.append(n_hits / n_relevant_total if n_relevant_total > 0 else 0.0)
            hits.append(1.0 if n_hits > 0 else 0.0)
            mean_sims.append(float(sim_row[top_k_idx].mean()))

    return MetricResult(
        precision=float(np.mean(precisions)),
        recall=float(np.mean(recalls)),
        hit_rate=float(np.mean(hits)),
        mean_cosine_sim=float(np.mean(mean_sims)),
        n_products=n,
    )


# ---------------------------------------------------------------------------
# 4. Run evaluation for all K values
# ---------------------------------------------------------------------------


def run_evaluation(
    embeddings: np.ndarray,
    product_ids: list[str],
    categories: list[str],
    k_values: tuple[int, ...] = DEFAULT_K_VALUES,
) -> dict[int, MetricResult]:
    results: dict[int, MetricResult] = {}
    for k in k_values:
        logger.info("Computing metrics @K=%d ...", k)
        t0 = time.perf_counter()
        results[k] = evaluate(embeddings, product_ids, categories, k)
        logger.info(
            "  @K=%2d  P=%.4f  R=%.4f  HR=%.4f  MeanCosSim=%.4f  (%.1fs)",
            k,
            results[k].precision,
            results[k].recall,
            results[k].hit_rate,
            results[k].mean_cosine_sim,
            time.perf_counter() - t0,
        )
    return results


# ---------------------------------------------------------------------------
# 5. Pretty-print results
# ---------------------------------------------------------------------------


def print_results(results: dict[int, MetricResult]) -> None:
    header = (
        f"{'K':>4}  {'Precision@K':>12}  {'Recall@K':>10}"
        f"  {'HitRate@K':>10}  {'MeanCosSim':>12}  {'Products':>10}"
    )
    sep = "=" * len(header)
    print(f"\n{sep}")
    print("Content-Based Filtering Evaluation (ecommerce_products_5k)")
    print(f"Embed model: {EMBED_MODEL}  |  Ground truth: same sub_category")
    print(sep)
    print(header)
    print("-" * len(header))
    for k in sorted(results):
        r = results[k]
        print(
            f"{k:>4}  {r.precision:>12.4f}  {r.recall:>10.4f}"
            f"  {r.hit_rate:>10.4f}  {r.mean_cosine_sim:>12.4f}  {r.n_products:>10}"
        )
    print(sep)


# ---------------------------------------------------------------------------
# 6. CLI
# ---------------------------------------------------------------------------


def _parse_args(argv: list[str] | None = None) -> argparse.Namespace:
    parser = argparse.ArgumentParser(
        description="Evaluate content-based filtering on ecommerce_products_5k.csv",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
    )
    parser.add_argument("--data-dir", type=Path, default=DEFAULT_DATA_DIR)
    parser.add_argument("--cache-dir", type=Path, default=DEFAULT_CACHE_DIR)
    parser.add_argument(
        "--k-values",
        type=lambda s: tuple(int(x) for x in s.split(",")),
        default=",".join(str(k) for k in DEFAULT_K_VALUES),
        help="Comma-separated cutoff values, e.g. 5,10,20",
    )
    parser.add_argument(
        "--no-cache",
        action="store_true",
        help="Re-embed even if a cache file exists (re-calls OpenAI API)",
    )
    parser.add_argument(
        "--env-file",
        type=Path,
        default=Path(".env"),
        help="Path to the .env file containing OPENAI_API_KEY",
    )
    parser.add_argument("--log-level", choices=["DEBUG", "INFO", "WARNING"], default="INFO")
    return parser.parse_args(argv)


def main(argv: list[str] | None = None) -> None:
    args = _parse_args(argv)
    logging.basicConfig(
        level=getattr(logging, args.log_level),
        format="%(asctime)s %(levelname)s %(name)s: %(message)s",
        stream=sys.stdout,
    )

    # Load .env
    env_path = args.env_file
    if env_path.exists():
        load_dotenv(env_path)
        logger.info("Loaded env from %s", env_path)
    else:
        logger.warning(".env file not found at %s — relying on environment variables", env_path)

    import os
    api_key = os.getenv("OPENAI_API_KEY")
    if not api_key:
        raise RuntimeError("OPENAI_API_KEY not set. Add it to your .env file or environment.")

    # Load products
    df = load_products(args.data_dir)
    categories = df[RELEVANCE_COL].tolist()

    # Embed (or load from cache)
    embeddings, product_ids = embed_products(
        df=df,
        api_key=api_key,
        cache_dir=args.cache_dir,
        use_cache=not args.no_cache,
    )

    # Evaluate
    logger.info("=== Starting evaluation ===")
    t0 = time.perf_counter()
    results = run_evaluation(embeddings, product_ids, categories, args.k_values)
    logger.info("Total evaluation time: %.1fs", time.perf_counter() - t0)

    print_results(results)


if __name__ == "__main__":
    main()

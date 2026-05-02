"""
Embed all active products from MongoDB into a Qdrant collection using
OpenAI text-embedding-ada-002 (vector size = 1536).

Usage:
    python -m scripts.embed_products          # from repo root
    python scripts/embed_products.py          # also works

What it does:
    1. Connects to MongoDB and reads the `products` collection.
    2. Resolves category names from the `categories` collection.
    3. Builds a rich text block per product (name, price, description,
       category, stock).
    4. Embeds the text with OpenAI `text-embedding-ada-002`.
    5. Creates / recreates the Qdrant collection `product_docs`
       (size=1536, COSINE) and upserts all points with full metadata
       in the payload.
"""

from __future__ import annotations

import logging
import os
import sys
import time
from pathlib import Path
from typing import Any

# ---------------------------------------------------------------------------
# Make sure `app.*` imports work when running as a standalone script.
# ---------------------------------------------------------------------------
_REPO_ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(_REPO_ROOT))

from dotenv import load_dotenv

load_dotenv(dotenv_path=_REPO_ROOT / ".env", override=False)

from openai import OpenAI
from pymongo import MongoClient
from qdrant_client import QdrantClient
from qdrant_client.http import models as qmodels

# ---------------------------------------------------------------------------
# Logging
# ---------------------------------------------------------------------------
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s %(levelname)s %(name)s: %(message)s",
)
logger = logging.getLogger("embed_products")

# ---------------------------------------------------------------------------
# Configuration — reads from .env (already loaded above)
# ---------------------------------------------------------------------------
MONGO_URI = os.environ["MONGO_URI"]
MONGO_DB_NAME = os.environ["MONGO_DB_NAME"]

QDRANT_URL = os.environ["QDRANT_URL"]
QDRANT_API_KEY = os.environ.get("QDRANT_API_KEY") or None
QDRANT_COLLECTION = os.environ.get("QDRANT_PRODUCT_DOC_COLLECTION", "product_docs")

OPENAI_API_KEY = os.environ["OPENAI_API_KEY"]
EMBEDDING_MODEL = "text-embedding-ada-002"  # 1536-dim
VECTOR_SIZE = 1536

BATCH_SIZE = 64  # OpenAI embedding API supports up to ~2048 inputs


# ---------------------------------------------------------------------------
# 1. MongoDB helpers
# ---------------------------------------------------------------------------

def _get_mongo_db():
    client = MongoClient(MONGO_URI)
    return client[MONGO_DB_NAME]


def _load_category_map(db) -> dict[str, str]:
    """Return {category_id: category_name} for all categories."""
    cat_map: dict[str, str] = {}
    for doc in db["categories"].find({}, {"_id": 1, "name": 1}):
        cat_map[str(doc["_id"])] = doc.get("name", "")
    logger.info("Loaded %d categories from MongoDB.", len(cat_map))
    return cat_map


def _load_products(db) -> list[dict[str, Any]]:
    """Return all active products."""
    query = {"is_active": True}
    products = list(db["products"].find(query))
    logger.info("Loaded %d active products from MongoDB.", len(products))
    return products


# ---------------------------------------------------------------------------
# 2. Text builder
# ---------------------------------------------------------------------------

def build_product_text(product: dict[str, Any], category_map: dict[str, str]) -> str:
    """
    Build a single rich text block for embedding.

    Includes: name, price range, description, categories, stock.
    """
    name = product.get("name", "")
    description = product.get("description", "")

    # Price range
    price_info = product.get("price", {})
    price_min = price_info.get("min")
    price_max = price_info.get("max")
    if price_min is not None and price_max is not None:
        if price_min == price_max:
            price_str = f"{price_min:,}đ"
        else:
            price_str = f"{price_min:,}đ – {price_max:,}đ"
    elif price_min is not None:
        price_str = f"{price_min:,}đ"
    else:
        price_str = "N/A"

    # Categories
    cat_ids = product.get("category_ids", [])
    cat_names = [category_map.get(cid, "") for cid in cat_ids]
    cat_names = [c for c in cat_names if c]  # filter empty
    category_str = ", ".join(cat_names) if cat_names else "N/A"

    # Stock
    stock = product.get("stock", 0)

    # Status
    status = product.get("status", "unknown")

    lines = [
        f"Tên sản phẩm: {name}",
        f"Giá: {price_str}",
        f"Danh mục: {category_str}",
        f"Tồn kho: {stock}",
        f"Trạng thái: {status}",
        f"Mô tả: {description}",
    ]
    return "\n".join(lines)


# ---------------------------------------------------------------------------
# 3. OpenAI embedding
# ---------------------------------------------------------------------------

def embed_texts(texts: list[str]) -> list[list[float]]:
    """Embed a list of texts using OpenAI text-embedding-ada-002."""
    client = OpenAI(api_key=OPENAI_API_KEY)
    all_embeddings: list[list[float]] = []

    for i in range(0, len(texts), BATCH_SIZE):
        batch = texts[i : i + BATCH_SIZE]
        logger.info(
            "  Embedding batch %d–%d of %d …",
            i + 1,
            min(i + BATCH_SIZE, len(texts)),
            len(texts),
        )
        response = client.embeddings.create(
            model=EMBEDDING_MODEL,
            input=batch,
        )
        # Response data is ordered by index
        batch_vectors = [item.embedding for item in response.data]
        all_embeddings.extend(batch_vectors)

    return all_embeddings


# ---------------------------------------------------------------------------
# 4. Qdrant helpers
# ---------------------------------------------------------------------------

def _get_qdrant_client() -> QdrantClient:
    return QdrantClient(url=QDRANT_URL, api_key=QDRANT_API_KEY)


def ensure_collection(client: QdrantClient) -> None:
    """Create the collection if it doesn't already exist."""
    collections = [c.name for c in client.get_collections().collections]

    if QDRANT_COLLECTION in collections:
        logger.info(
            "Collection '%s' already exists — it will be recreated.",
            QDRANT_COLLECTION,
        )
        client.delete_collection(QDRANT_COLLECTION)

    client.create_collection(
        collection_name=QDRANT_COLLECTION,
        vectors_config=qmodels.VectorParams(
            size=VECTOR_SIZE,
            distance=qmodels.Distance.COSINE,
        ),
    )
    logger.info(
        "Created Qdrant collection '%s' (size=%d, COSINE).",
        QDRANT_COLLECTION,
        VECTOR_SIZE,
    )


def build_payload(product: dict[str, Any], category_map: dict[str, str]) -> dict[str, Any]:
    """Build the Qdrant payload with full product metadata."""
    cat_ids = product.get("category_ids", [])
    cat_names = [category_map.get(cid, "") for cid in cat_ids if category_map.get(cid)]

    price_info = product.get("price", {})

    return {
        "product_id": str(product["_id"]),
        "name": product.get("name", ""),
        "description": product.get("description", ""),
        "status": product.get("status", ""),
        "is_active": product.get("is_active", False),
        "price_min": price_info.get("min"),
        "price_max": price_info.get("max"),
        "stock": product.get("stock", 0),
        "rating": product.get("rating", 0.0),
        "rate_count": product.get("rate_count", 0),
        "sold_count": product.get("sold_count", 0),
        "seller_id": product.get("seller_id", ""),
        "category_ids": cat_ids,
        "category_names": cat_names,
    }


def upsert_to_qdrant(
    client: QdrantClient,
    products: list[dict[str, Any]],
    embeddings: list[list[float]],
    category_map: dict[str, str],
) -> int:
    """Upsert product embeddings to Qdrant. Returns the count of upserted points."""
    UPSERT_BATCH = 100
    total = 0

    for i in range(0, len(products), UPSERT_BATCH):
        batch_products = products[i : i + UPSERT_BATCH]
        batch_vectors = embeddings[i : i + UPSERT_BATCH]

        points = []
        for product, vector in zip(batch_products, batch_vectors):
            point = qmodels.PointStruct(
                id=str(product["_id"]),
                vector=vector,
                payload=build_payload(product, category_map),
            )
            points.append(point)

        client.upsert(collection_name=QDRANT_COLLECTION, points=points)
        total += len(points)
        logger.info("  Upserted %d / %d points.", total, len(products))

    return total


# ---------------------------------------------------------------------------
# Main
# ---------------------------------------------------------------------------

def main() -> None:
    logger.info("=== embed_products started ===")
    t0 = time.perf_counter()

    # 1. Load data from MongoDB
    db = _get_mongo_db()
    category_map = _load_category_map(db)
    products = _load_products(db)

    if not products:
        logger.warning("No active products found. Exiting.")
        return

    # 2. Build text blocks
    logger.info("Building text blocks for %d products …", len(products))
    texts = [build_product_text(p, category_map) for p in products]

    # Log a sample
    logger.info("--- Sample text block (first product) ---\n%s\n---", texts[0])

    # 3. Embed with OpenAI
    logger.info("Embedding %d texts with %s …", len(texts), EMBEDDING_MODEL)
    embeddings = embed_texts(texts)
    assert len(embeddings) == len(products), "Embedding count mismatch!"

    # 4. Upsert to Qdrant
    logger.info("Upserting to Qdrant collection '%s' …", QDRANT_COLLECTION)
    qdrant = _get_qdrant_client()
    ensure_collection(qdrant)
    count = upsert_to_qdrant(qdrant, products, embeddings, category_map)

    elapsed = time.perf_counter() - t0
    logger.info(
        "=== embed_products completed in %.2fs. %d products embedded & stored. ===",
        elapsed,
        count,
    )


if __name__ == "__main__":
    main()

"""
MongoDB persistence for user–product interaction history (ai-service).
"""

from __future__ import annotations

import logging
from datetime import datetime, timezone
from functools import lru_cache
from typing import Any

from pymongo import ASCENDING, DESCENDING, MongoClient
from pymongo.collection import Collection

from app.core.config import get_settings

logger = logging.getLogger(__name__)

_indexes_ensured = False


@lru_cache
def get_mongo_client() -> MongoClient:
    settings = get_settings()
    if not settings.MONGO_URI:
        raise RuntimeError("MONGO_URI is not set")
    return MongoClient(settings.MONGO_URI)


def get_interactions_collection() -> Collection:
    settings = get_settings()
    if not settings.MONGO_DB_NAME:
        raise RuntimeError("MONGO_DB_NAME is not set")
    db = get_mongo_client()[settings.MONGO_DB_NAME]
    return db[settings.INTERACTION_COLLECTION_NAME]


def ensure_interaction_indexes() -> None:
    """Idempotent: create indexes for user history queries."""
    global _indexes_ensured
    if _indexes_ensured:
        return
    col = get_interactions_collection()
    col.create_index([("user_id", ASCENDING), ("timestamp", DESCENDING)])
    col.create_index([("user_id", ASCENDING), ("product_id", ASCENDING)])
    _indexes_ensured = True
    logger.info("Ensured indexes on %s", col.name)


def insert_interaction(
    *,
    user_id: str,
    product_id: str,
    action: str,
    score: float,
    timestamp: datetime | None = None,
) -> None:
    col = get_interactions_collection()
    ts = timestamp if timestamp is not None else datetime.now(timezone.utc)
    if ts.tzinfo is None:
        ts = ts.replace(tzinfo=timezone.utc)

    doc: dict[str, Any] = {
        "user_id": user_id,
        "product_id": product_id,
        "action": action,
        "score": float(score),
        "timestamp": ts,
    }
    col.insert_one(doc)


def get_recent_interactions_for_user(user_id: str, limit: int) -> list[dict[str, Any]]:
    """
    Return the newest `limit` interactions for a user (newest first).
    """
    col = get_interactions_collection()
    cursor = (
        col.find({"user_id": user_id}).sort("timestamp", DESCENDING).limit(limit)
    )
    return list(cursor)

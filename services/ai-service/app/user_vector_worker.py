"""
RabbitMQ worker: persist user–product interactions to MongoDB, then refresh the
user vector in Qdrant from the newest N stored interactions.

Primary message format (JSON):
{
  "user_id": "u123",
  "product_id": "p456",
  "action": "view",
  "score": 1,
  "timestamp": "2026-03-22T20:00:00Z"
}

`score` is optional if `INTERACTION_ACTION_SCORES_JSON` defines weights per `action`.
`timestamp` is optional (defaults to server UTC now).

Legacy format (still supported):
{
  "user_id": "USER_UUID",
  "items": [
    { "product_id": "PRODUCT_UUID_1", "weight": 1.0 }
  ]
}
"""

from __future__ import annotations

import json
import logging
from typing import Any

import pika

from app.core.config import get_settings
from app.schemas.user_interaction import UserInteractionMessage
from app.schemas.user_vector import UserVectorRequest
from app.services.interaction_service import (
    ensure_interaction_indexes,
    get_recent_interactions_for_user,
    insert_interaction,
)
from app.services.user_vector_service import (
    compute_user_vector,
    compute_user_vector_from_interaction_docs,
    upsert_user_vector,
)

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)


def _resolve_interaction_score(action: str, explicit_score: float | None) -> float:
    settings = get_settings()
    if explicit_score is not None:
        return float(explicit_score)
    return float(
        settings.INTERACTION_ACTION_SCORES.get(
            action, settings.INTERACTION_DEFAULT_SCORE
        )
    )


def _handle_legacy_user_vector_request(data: dict[str, Any]) -> None:
    req = UserVectorRequest(**data)
    if not req.items:
        raise ValueError("items must not be empty")
    vector = compute_user_vector(req.items)
    upsert_user_vector(req.user_id, vector)
    logger.info("Updated user vector for %s (legacy items payload)", req.user_id)


def _handle_interaction_message(data: dict[str, Any]) -> None:
    settings = get_settings()
    msg = UserInteractionMessage.model_validate(data)

    score = _resolve_interaction_score(msg.action, msg.score)

    insert_interaction(
        user_id=msg.user_id,
        product_id=msg.product_id,
        action=msg.action,
        score=score,
        timestamp=msg.timestamp,
    )

    limit = settings.USER_INTERACTION_HISTORY_LIMIT
    recent = get_recent_interactions_for_user(msg.user_id, limit)
    if not recent:
        raise RuntimeError("Interaction inserted but no rows returned for user (unexpected)")

    vector = compute_user_vector_from_interaction_docs(recent)
    upsert_user_vector(msg.user_id, vector)

    logger.info(
        "Stored interaction and updated user vector for %s (used %s recent interactions, limit=%s)",
        msg.user_id,
        len(recent),
        limit,
    )


def _handle_message(body: bytes) -> None:
    settings = get_settings()
    if not settings.MONGO_URI or not settings.MONGO_DB_NAME:
        raise RuntimeError(
            "user_vector_worker requires MONGO_URI and MONGO_DB_NAME to be set"
        )

    data: dict[str, Any] = json.loads(body.decode("utf-8"))

    # Legacy: pre-aggregated product weights (no Mongo write)
    if isinstance(data.get("items"), list):
        _handle_legacy_user_vector_request(data)
        return

    _handle_interaction_message(data)


def main() -> None:
    settings = get_settings()

    if not settings.MONGO_URI or not settings.MONGO_DB_NAME:
        raise RuntimeError(
            "Set MONGO_URI and MONGO_DB_NAME for user_vector_worker (MongoDB interaction store)."
        )

    ensure_interaction_indexes()

    params = pika.URLParameters(settings.RABBITMQ_URL)
    connection = pika.BlockingConnection(params)
    channel = connection.channel()

    queue_name = settings.USER_VECTOR_QUEUE
    channel.queue_declare(queue=queue_name, durable=True)

    def callback(ch, method, properties, body) -> None:  # type: ignore[no-untyped-def]
        try:
            _handle_message(body)
            ch.basic_ack(delivery_tag=method.delivery_tag)
        except Exception:
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

    channel.basic_qos(prefetch_count=1)
    channel.basic_consume(queue=queue_name, on_message_callback=callback)

    logger.info(
        "User vector worker started. Listening on queue '%s' at %s (history_limit=%s)",
        queue_name,
        settings.RABBITMQ_URL,
        settings.USER_INTERACTION_HISTORY_LIMIT,
    )

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        logger.info("Stopping user vector worker...")
        channel.stop_consuming()
    finally:
        connection.close()


if __name__ == "__main__":
    main()

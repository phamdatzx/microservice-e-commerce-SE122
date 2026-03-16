"""
RabbitMQ worker for computing and storing user vectors in Qdrant.

Expected message format (JSON), compatible with UserVectorRequest:
{
  "user_id": "USER_UUID",
  "items": [
    { "product_id": "PRODUCT_UUID_1", "weight": 1.0 },
    { "product_id": "PRODUCT_UUID_2", "weight": 0.5 }
  ]
}
"""

import json
import logging
from typing import Any

import pika

from app.core.config import get_settings
from app.schemas.user_vector import UserVectorRequest
from app.services.user_vector_service import compute_user_vector, upsert_user_vector

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)


def _handle_message(body: bytes) -> None:
    try:
        data: dict[str, Any] = json.loads(body.decode("utf-8"))
        req = UserVectorRequest(**data)

        if not req.items:
            raise ValueError("items must not be empty")

        vector = compute_user_vector(req.items)
        upsert_user_vector(req.user_id, vector)

        logger.info("Updated user vector for %s", req.user_id)
    except Exception:  # noqa: BLE001
        logger.exception("Failed to process user vector message: %s", body)
        raise


def main() -> None:
    settings = get_settings()

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
        "User vector worker started. Listening on queue '%s' at %s",
        queue_name,
        settings.RABBITMQ_URL,
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


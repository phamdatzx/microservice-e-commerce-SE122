"""
Basic RabbitMQ worker for indexing products into Qdrant.

Expected message format (JSON), compatible with ProductPayload:
{
  "id": "product-uuid",
  "name": "iPhone 15 Pro Max",
  "category_id": "1",
  "category_name": "phone",
  "seller_id": "seller-uuid"
}
"""

import json
import logging
from typing import Any

import pika

from app.core.config import get_settings
from app.schemas.product_vector import ProductPayload, ProductVector
from app.services.embedding_service import compute_product_embedding
from app.services.qdrant_service import upsert_product_vector

logger = logging.getLogger(__name__)
logging.basicConfig(level=logging.INFO)


def _handle_message(body: bytes) -> None:
    settings = get_settings()

    try:
        data: dict[str, Any] = json.loads(body.decode("utf-8"))
        payload = ProductPayload(**data)
        vector = compute_product_embedding(payload)
        product_vec = ProductVector(id=payload.id, vector=vector, payload=payload)

        upsert_product_vector(product_vec)
        logger.info("Indexed product %s into Qdrant", payload.id)
    except Exception:  # noqa: BLE001
        logger.exception("Failed to process product message: %s", body)
        # Let the exception bubble up to control ack/nack policy at call site.
        raise


def main() -> None:
    settings = get_settings()

    params = pika.URLParameters(settings.RABBITMQ_URL)
    connection = pika.BlockingConnection(params)
    channel = connection.channel()

    # Declare queue to ensure it exists (idempotent).
    channel.queue_declare(queue=settings.PRODUCT_CREATED_QUEUE, durable=True)

    def callback(ch, method, properties, body) -> None:  # type: ignore[no-untyped-def]
        try:
            _handle_message(body)
            ch.basic_ack(delivery_tag=method.delivery_tag)
        except Exception:
            # For now, nack without requeue to avoid poison message loops.
            ch.basic_nack(delivery_tag=method.delivery_tag, requeue=False)

    channel.basic_qos(prefetch_count=1)
    channel.basic_consume(
        queue=settings.PRODUCT_CREATED_QUEUE,
        on_message_callback=callback,
    )

    logger.info(
        "Worker started. Listening on queue '%s' at %s",
        settings.PRODUCT_CREATED_QUEUE,
        settings.RABBITMQ_URL,
    )

    try:
        channel.start_consuming()
    except KeyboardInterrupt:
        logger.info("Stopping worker...")
        channel.stop_consuming()
    finally:
        connection.close()


if __name__ == "__main__":
    main()


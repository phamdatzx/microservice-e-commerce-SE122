import json
import pika

RABBITMQ_URL = "amqp://guest:guest@localhost:5672/"  # same as RABBITMQ_URL
QUEUE = "user.vector.update"  # same as USER_VECTOR_QUEUE

body = {
    "user_id": "550e8400-e29b-41d4-a716-446655440000",
    "product_id": "01490549-83a2-450b-884e-91629a9978b5",
    "action": "purchase",
    # "score": 1,  # optional if you use INTERACTION_ACTION_SCORES_JSON
    # "timestamp": "2026-03-22T20:00:00Z",  # optional
}

params = pika.URLParameters(RABBITMQ_URL)
connection = pika.BlockingConnection(params)
channel = connection.channel()

channel.queue_declare(queue=QUEUE, durable=True)
channel.basic_publish(
    exchange="",
    routing_key=QUEUE,
    body=json.dumps(body).encode("utf-8"),
    properties=pika.BasicProperties(
        content_type="application/json",
        delivery_mode=2,  # persistent (if broker supports it)
    ),
)
connection.close()
print("Sent:", body)
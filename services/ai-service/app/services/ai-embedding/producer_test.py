import json
import uuid
import time
import pika
from config import RABBITMQ_URL, RABBITMQ_QUEUE

def send_test(n=3):
    params = pika.URLParameters(RABBITMQ_URL)
    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(queue=RABBITMQ_QUEUE, durable=True)

    # 1. Test product messages
    for i in range(n):
        product = {
            "id": str(uuid.uuid4()),
            "name": f"Test Product {i}",
            "description": f"Mô tả sản phẩm thử nghiệm số {i}, có tính năng đặc biệt.",
            "category_ids": ["cat-123"],
            "seller_id": "seller-xyz",
            "price": {"min": 1000000 + i*50000, "max": 1000000 + i*50000},
            "rating": 4,
            "rate_count": 10,
            "sold_count": 5,
            "stock": 20,
            "is_active": True,
            "updated_at": "2026-03-30T22:30:00Z"
        }
        payload = {"type": "product", "payload": product}
        ch.basic_publish(exchange="", routing_key=RABBITMQ_QUEUE, body=json.dumps(payload))
        print("Sent product", product["id"])
        time.sleep(0.1)

    # 2. Test seller messages
    for i in range(n):
        seller = {
            "id": str(uuid.uuid4()),
            "name": f"Test Seller {i}",
            "address": {
                "ward": "Phường Thảo Điền",
                "district": "Thành Phố Thủ Đức",
                "province": "Hồ Chí Minh",
                "country": "Việt Nam"
            },
            "sale_info": {
                "follow_count": i,
                "rating_count": 10+i,
                "rating_average": 4.5,
                "product_count": 20+i
            }
        }
        payload = {"type": "seller", "payload": seller}
        ch.basic_publish(exchange="", routing_key=RABBITMQ_QUEUE, body=json.dumps(payload))
        print("Sent seller", seller["id"])
        time.sleep(0.1)

    # 3. Test rating messages
    for i in range(n):
        rating = {
            "id": str(uuid.uuid4()),
            "product_id": "prod-abc",
            "variant_id": "var-xyz",
            "user": {"id": "user-123", "name": "tester"},
            "star": 5,
            "content": f"Sản phẩm rất tốt, đánh giá thử nghiệm số {i}.",
            "created_at": "2026-03-30T22:30:00Z"
        }
        payload = {"type": "rating", "payload": rating}
        ch.basic_publish(exchange="", routing_key=RABBITMQ_QUEUE, body=json.dumps(payload))
        print("Sent rating", rating["id"])
        time.sleep(0.1)

    conn.close()

if __name__ == "__main__":
    send_test(3)

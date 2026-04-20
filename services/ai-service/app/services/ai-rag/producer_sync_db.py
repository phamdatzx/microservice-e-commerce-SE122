import os
import json
import pika
from pymongo import MongoClient
from config import RABBITMQ_URL, RABBITMQ_QUEUE

# Mình đã lấy URI thật từ file product-service/.env của bạn để tránh lỗi Connection Refused trên local.
MONGO_URI = os.getenv("MONGO_URI", "mongodb+srv://megumikatoubuhbuh_db_user:L3tM31n!bro@e-commerce.xokv4hb.mongodb.net/?appName=e-commerce")
DB_NAME = os.getenv("DB_NAME", "ecommerce")

# ---------------------------------------------------------
# 1. Database connections
# ---------------------------------------------------------
mongo_client = MongoClient(MONGO_URI)
db = mongo_client[DB_NAME]
product_collection = db["products"]
rating_collection = db["ratings"]

# ---------------------------------------------------------
# 2. RabbitMQ helper
# ---------------------------------------------------------
def send_to_queue(msg):
    params = pika.URLParameters(RABBITMQ_URL)
    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(queue=RABBITMQ_QUEUE, durable=True)
    ch.basic_publish(exchange="", routing_key=RABBITMQ_QUEUE, body=json.dumps(msg))
    conn.close()

# ---------------------------------------------------------
# 3. Sync Products
# ---------------------------------------------------------
def sync_products():
    print("Syncing products from MongoDB Atlas...")
    count = 0
    for doc in product_collection.find():
        product = {
            "id": str(doc.get("_id", "")),
            "name": doc.get("name", ""),
            "description": doc.get("description", ""),
            "seller_id": doc.get("seller_id", ""),
            "price_min": doc.get("price", {}).get("min", 0),
            "price_max": doc.get("price", {}).get("max", 0),
            "rating": doc.get("rating", 0.0),
            "rate_count": doc.get("rate_count", 0),
            "sold_count": doc.get("sold_count", 0),
            "stock": doc.get("stock", 0),
            "is_active": doc.get("is_active", False),
            "category_ids": doc.get("category_ids", []),
            "updated_at": doc.get("updated_at").isoformat() if hasattr(doc.get("updated_at"), "isoformat") else doc.get("updated_at")
        }

        msg = {
            "type": "product",
            "payload": product
        }

        send_to_queue(msg)
        count += 1
        print("Sent product:", product["name"])
        
    print(f"=> Done syncing {count} products!")

# ---------------------------------------------------------
# 4. Sync Ratings
# ---------------------------------------------------------
def sync_ratings():
    print("Syncing ratings from MongoDB Atlas...")
    count = 0
    for doc in rating_collection.find():
        rating = {
            "id": str(doc.get("_id", "")),
            "product_id": doc.get("product_id", ""),
            "variant_id": doc.get("variant_id", ""),
            "user": {
                "id": str(doc.get("user", {}).get("_id", "")),
                "name": doc.get("user", {}).get("name", "")
            },
            "star": doc.get("star", 0),
            "content": doc.get("content", ""),
            "created_at": doc.get("created_at").isoformat() if hasattr(doc.get("created_at"), "isoformat") else doc.get("created_at")
        }

        msg = {
            "type": "rating",
            "payload": rating
        }

        send_to_queue(msg)
        count += 1
        print("Sent rating for product_id:", rating["product_id"])

    print(f"=> Done syncing {count} ratings!")

# ---------------------------------------------------------
# 5. Run all
# ---------------------------------------------------------
if __name__ == "__main__":
    # Theo yêu cầu của user: Bỏ qua Sync Seller, chỉ tập trung Product & Rating
    sync_products()
    sync_ratings()

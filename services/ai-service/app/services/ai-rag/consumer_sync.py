import json
import logging
import threading
import time
from typing import List

import pika
from qdrant_client import QdrantClient
from qdrant_client.models import VectorParams, Distance

from config import (
    RABBITMQ_URL, RABBITMQ_QUEUE, RABBITMQ_DLQ,
    QDRANT_HOST, QDRANT_PORT, 
    QDRANT_PRODUCT_COLLECTION, QDRANT_RATING_COLLECTION,
    EMBEDDING_PROVIDER, EMBEDDING_MODEL, OPENAI_API_KEY, BATCH_SIZE, FLUSH_INTERVAL
)

logging.basicConfig(level=logging.INFO)
logger = logging.getLogger("ai-embedding-consumer-sync")

# Embedding provider abstraction (sync)
class EmbeddingModel:
    def __init__(self):
        if EMBEDDING_PROVIDER == "sentence_transformers":
            from sentence_transformers import SentenceTransformer
            self.model = SentenceTransformer(EMBEDDING_MODEL)
            self.dim = self.model.get_sentence_embedding_dimension()
            self._encode = self._st_encode
        elif EMBEDDING_PROVIDER == "openai":
            import openai
            openai.api_key = OPENAI_API_KEY
            self.dim = 1536
            self._encode = self._openai_encode
        else:
            raise ValueError("Unsupported EMBEDDING_PROVIDER")

    def _st_encode(self, texts: List[str]) -> List[List[float]]:
        return self.model.encode(texts, show_progress_bar=False).tolist()

    def _openai_encode(self, texts: List[str]) -> List[List[float]]:
        import openai
        resp = openai.Embedding.create(model="text-embedding-3-small", input=texts)
        return [e["embedding"] for e in resp["data"]]

    def encode(self, texts: List[str]) -> List[List[float]]:
        return self._encode(texts)

# Qdrant client
qdrant = QdrantClient(host=QDRANT_HOST, port=QDRANT_PORT)

def ensure_collections(dim):
    collections = [QDRANT_PRODUCT_COLLECTION, QDRANT_RATING_COLLECTION]
    for collection in collections:
        try:
            try:
                col = qdrant.get_collection(collection_name=collection)
                logger.info("Qdrant collection %s exists", collection)
            except Exception:
                qdrant.recreate_collection(
                    collection_name=collection,
                    vectors_config=VectorParams(size=dim, distance=Distance.COSINE)
                )
                logger.info("Created Qdrant collection %s dim=%s", collection, dim)
        except Exception as e:
            logger.exception("Error ensuring collection %s: %s", collection, e)
            raise

def build_point(msg_type, obj, vector):
    if msg_type == "product":
        payload = {
            "product_id": obj["id"],
            "title": obj.get("name"),
            "category_ids": obj.get("category_ids", []),
            "seller_id": obj.get("seller_id"),
            "price_min": obj.get("price_min"),  # Sửa lại tên field từ dict
            "price_max": obj.get("price_max"),
            "rating": obj.get("rating"),
            "rate_count": obj.get("rate_count"),
            "sold_count": obj.get("sold_count"),
            "stock": obj.get("stock"),
            "is_active": obj.get("is_active"),
            "updated_at": obj.get("updated_at"),
            "embedding_version": "v1",
        }
        return {"id": obj["id"], "vector": vector, "payload": payload}

    elif msg_type == "rating":
        payload = {
            "rating_id": obj["id"],
            "product_id": obj.get("product_id"),
            "variant_id": obj.get("variant_id"),
            "user_id": obj.get("user", {}).get("id"),
            "user_name": obj.get("user", {}).get("name"),
            "star": obj.get("star"),
            "content": obj.get("content"),
            "created_at": obj.get("created_at"),
            "embedding_version": "v1",
        }
        return {"id": obj["id"], "vector": vector, "payload": payload}

    return None

class SyncEmbeddingConsumer:
    def __init__(self):
        self.model = EmbeddingModel()
        ensure_collections(self.model.dim)
        self.buffer = []
        self.lock = threading.Lock()
        self._stop = threading.Event()
        self.flush_thread = threading.Thread(target=self._periodic_flush, daemon=True)
        self.flush_thread.start()

    def _periodic_flush(self):
        while not self._stop.is_set():
            time.sleep(FLUSH_INTERVAL)
            self.flush_buffer()

    def stop(self):
        self._stop.set()
        self.flush_thread.join(timeout=2)

    def on_message(self, ch, method, properties, body):
        try:
            payload = json.loads(body.decode())
            msg_type = payload.get("type")
            obj = payload.get("payload", {})
            
            # Chỉ xử lý product và rating, bỏ qua event ko hợp lệ / out_of_scope
            if not msg_type or not obj or msg_type not in ["product", "rating"]:
                ch.basic_ack(delivery_tag=method.delivery_tag)
                return

            with self.lock:
                self.buffer.append({"type": msg_type, "payload": obj})
                if len(self.buffer) >= BATCH_SIZE:
                    self.flush_buffer()

            ch.basic_ack(delivery_tag=method.delivery_tag)
        except Exception:
            logger.exception("Processing failed, sending to DLQ")
            self.send_to_dlq(body)
            ch.basic_ack(delivery_tag=method.delivery_tag)

    def flush_buffer(self):
        with self.lock:
            if not self.buffer:
                return
            batch = self.buffer
            self.buffer = []

        texts, objs, types = [], [], []
        for b in batch:
            msg_type = b["type"]
            obj = b["payload"]
            
            # Format đoạn text chuẩn đưa vào Vector theo chiến lược
            if msg_type == "product":
                text = f"Tên sản phẩm: {obj.get('name','')} - Mô tả và thông số kỹ thuật: {obj.get('description','')}"
            elif msg_type == "rating":
                text = f"Đánh giá sản phẩm: {obj.get('content','')}"
            else:
                continue
                
            texts.append(text)
            objs.append(obj)
            types.append(msg_type)

        if not texts:
            return

        try:
            vectors = self.model.encode(texts)
        except Exception:
            logger.exception("Embedding generation failed; requeueing batch to DLQ")
            for b in batch:
                self.send_to_dlq(json.dumps(b).encode())
            return

        points_by_collection = {
            QDRANT_PRODUCT_COLLECTION: [],
            QDRANT_RATING_COLLECTION: []
        }

        for t, o, v in zip(types, objs, vectors):
            point = build_point(t, o, v)
            if not point: continue
            
            if t == "product":
                points_by_collection[QDRANT_PRODUCT_COLLECTION].append(point)
            elif t == "rating":
                points_by_collection[QDRANT_RATING_COLLECTION].append(point)

        for coll, points in points_by_collection.items():
            if not points: continue
            try:
                qdrant.upsert(collection_name=coll, points=points)
                logger.info("Upserted %d points to Qdrant collection %s", len(points), coll)
            except Exception:
                logger.exception("Qdrant upsert failed for %s, retrying once", coll)
                try:
                    qdrant.upsert(collection_name=coll, points=points)
                except Exception:
                    logger.exception("Retry failed, sending points to DLQ")
                    for p in points:
                        # Re-send payload
                        self.send_to_dlq(json.dumps({"id": p["id"], "type": "error", "payload": p.get("payload", {})}).encode())

    def send_to_dlq(self, body: bytes):
        try:
            params = pika.URLParameters(RABBITMQ_URL)
            conn = pika.BlockingConnection(params)
            ch = conn.channel()
            ch.queue_declare(queue=RABBITMQ_DLQ, durable=True)
            ch.basic_publish(exchange="", routing_key=RABBITMQ_DLQ, body=body)
            conn.close()
        except Exception:
            logger.exception("Failed to publish to DLQ")

def main():
    consumer = SyncEmbeddingConsumer()
    params = pika.URLParameters(RABBITMQ_URL)
    conn = pika.BlockingConnection(params)
    ch = conn.channel()
    ch.queue_declare(queue=RABBITMQ_QUEUE, durable=True)
    ch.queue_declare(queue=RABBITMQ_DLQ, durable=True)
    logger.info("Waiting for messages on %s", RABBITMQ_QUEUE)
    ch.basic_qos(prefetch_count=BATCH_SIZE * 2)
    ch.basic_consume(queue=RABBITMQ_QUEUE, on_message_callback=consumer.on_message)
    try:
        ch.start_consuming()
    except KeyboardInterrupt:
        logger.info("Stopping consumer...")
        ch.stop_consuming()
        consumer.stop()
        conn.close()

if __name__ == "__main__":
    main()

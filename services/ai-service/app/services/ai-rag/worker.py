import pika
import json
from sentence_transformers import SentenceTransformer

model = SentenceTransformer("BAAI/bge-m3")

def build_text(product):
    return f"""
    passage:
    Tên sản phẩm: {product['name']}.
    Danh mục: {product['category']}.
    Giá: {product['price']} VND.
    Mô tả: {product['description']}.
    """

def callback(ch, method, properties, body):
    product = json.loads(body)

    text = build_text(product)
    embedding = model.encode(text)

    # 👉 save to FAISS / Qdrant
    print("Embedded:", product["id"])

    ch.basic_ack(delivery_tag=method.delivery_tag)

connection = pika.BlockingConnection(
    pika.ConnectionParameters("localhost")
)
channel = connection.channel()

channel.queue_declare(queue="embedding_queue")

channel.basic_consume(
    queue="embedding_queue",
    on_message_callback=callback
)

print("Waiting for messages...")
channel.start_consuming()
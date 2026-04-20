from qdrant_client import QdrantClient
qc = QdrantClient(host="localhost", port=6333)

from config import (
    RABBITMQ_URL, RABBITMQ_QUEUE, RABBITMQ_DLQ,
    QDRANT_HOST, QDRANT_PORT, 
    QDRANT_PRODUCT_COLLECTION, QDRANT_RATING_COLLECTION,
    EMBEDDING_PROVIDER, EMBEDDING_MODEL, OPENAI_API_KEY, BATCH_SIZE, FLUSH_INTERVAL
)

print(qc.get_collection(QDRANT_PRODUCT_COLLECTION))
print(qc.scroll(collection_name=QDRANT_PRODUCT_COLLECTION, limit=10))
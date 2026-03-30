import os

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
RABBITMQ_QUEUE = os.getenv("RABBITMQ_QUEUE", "ai.embedding.requests")
RABBITMQ_DLQ = os.getenv("RABBITMQ_DLQ", "ai.embedding.dlq")

QDRANT_HOST = os.getenv("QDRANT_HOST", "localhost")
QDRANT_PORT = int(os.getenv("QDRANT_PORT", 6333))
QDRANT_COLLECTION = os.getenv("QDRANT_COLLECTION", "docs")

EMBEDDING_PROVIDER = os.getenv("EMBEDDING_PROVIDER", "sentence_transformers")
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY", "")

BATCH_SIZE = int(os.getenv("BATCH_SIZE", 16))
FLUSH_INTERVAL = int(os.getenv("FLUSH_INTERVAL", 5))  # seconds
MAX_RETRY = int(os.getenv("MAX_RETRY", 3))

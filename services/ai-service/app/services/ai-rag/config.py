import os

RABBITMQ_URL = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
RABBITMQ_QUEUE = os.getenv("RABBITMQ_QUEUE", "ai.embedding.requests")
RABBITMQ_DLQ = os.getenv("RABBITMQ_DLQ", "ai.embedding.dlq")

QDRANT_HOST = os.getenv("QDRANT_HOST", "localhost")
QDRANT_PORT = int(os.getenv("QDRANT_PORT", 6333))

# Chia tách 2 collection rõ ràng
QDRANT_PRODUCT_COLLECTION = os.getenv("QDRANT_PRODUCT_COLLECTION", "products")
QDRANT_RATING_COLLECTION = os.getenv("QDRANT_RATING_COLLECTION", "ratings")

EMBEDDING_PROVIDER = os.getenv("EMBEDDING_PROVIDER", "sentence_transformers")
EMBEDDING_MODEL = os.getenv("EMBEDDING_MODEL", "BAAI/bge-m3")

# LLM Configuration (Llama 3 via OpenRouter)
OPENROUTER_API_KEY = os.getenv("OPENROUTER_API_KEY", "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")
OPENROUTER_BASE_URL = os.getenv("OPENROUTER_BASE_URL", "https://openrouter.ai/api/v1")
LLM_MODEL = os.getenv("LLM_MODEL", "meta-llama/llama-3-8b-instruct:free")

BATCH_SIZE = int(os.getenv("BATCH_SIZE", 16))
FLUSH_INTERVAL = int(os.getenv("FLUSH_INTERVAL", 5))  # seconds
MAX_RETRY = int(os.getenv("MAX_RETRY", 3))

import os
from functools import lru_cache


class Settings:
    """Service settings loaded from environment variables."""

    PROJECT_NAME: str = os.getenv("PROJECT_NAME", "AI Service")
    API_V1_PREFIX: str = os.getenv("API_V1_PREFIX", "")
    VERSION: str = os.getenv("VERSION", "0.1.0")
    DEBUG: bool = os.getenv("DEBUG", "false").lower() == "true"

    # Qdrant configuration
    QDRANT_URL: str = os.getenv("QDRANT_URL", "http://localhost:6333")
    QDRANT_API_KEY: str | None = os.getenv("QDRANT_API_KEY")
    QDRANT_COLLECTION_NAME: str = os.getenv("QDRANT_COLLECTION_NAME", "products")
    QDRANT_USER_COLLECTION_NAME: str = os.getenv(
        "QDRANT_USER_COLLECTION_NAME",
        "user_vectors",
    )

    # RabbitMQ configuration (for background workers)
    RABBITMQ_URL: str = os.getenv("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/")
    PRODUCT_CREATED_QUEUE: str = os.getenv(
        "PRODUCT_CREATED_QUEUE",
        "product.created",
    )
    USER_VECTOR_QUEUE: str = os.getenv(
        "USER_VECTOR_QUEUE",
        "user.vector.update",
    )

    # Embedding model configuration
    # Default to a lighter, fast model; override via EMBEDDING_MODEL_NAME if needed.
    EMBEDDING_MODEL_NAME: str = os.getenv(
        "EMBEDDING_MODEL_NAME",
        "intfloat/multilingual-e5-small",
    )


@lru_cache
def get_settings() -> Settings:
    return Settings()

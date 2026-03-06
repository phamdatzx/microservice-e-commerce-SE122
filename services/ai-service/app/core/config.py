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

    # Embedding model configuration
    EMBEDDING_MODEL_NAME: str = os.getenv("EMBEDDING_MODEL_NAME", "BAAI/bge-m3")


@lru_cache
def get_settings() -> Settings:
    return Settings()

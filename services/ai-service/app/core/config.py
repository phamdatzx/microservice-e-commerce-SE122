import os
from pathlib import Path
from functools import lru_cache

from dotenv import load_dotenv


_ENV_PATH = Path(__file__).resolve().parents[2] / ".env"
# Load .env if present, but never override real environment variables.
load_dotenv(dotenv_path=_ENV_PATH, override=False)


def _require_env(name: str) -> str:
    value = os.getenv(name)
    if value is None:
        raise RuntimeError(
            f"Missing required environment variable: {name}. "
            f"Set it in the environment or in '{_ENV_PATH}'."
        )
    return value


def _env_bool(name: str) -> bool:
    raw = _require_env(name).strip().lower()
    if raw in {"true", "1", "yes", "y", "on"}:
        return True
    if raw in {"false", "0", "no", "n", "off"}:
        return False
    raise RuntimeError(
        f"Invalid boolean value for {name}: {raw!r}. "
        "Use one of: true/false, 1/0, yes/no, on/off."
    )


class Settings:
    """Service settings loaded from environment variables."""

    PROJECT_NAME: str = _require_env("PROJECT_NAME")
    API_V1_PREFIX: str = _require_env("API_V1_PREFIX")
    VERSION: str = _require_env("VERSION")
    DEBUG: bool = _env_bool("DEBUG")

    # Qdrant configuration
    QDRANT_URL: str = _require_env("QDRANT_URL")
    QDRANT_API_KEY: str | None = os.getenv("QDRANT_API_KEY") or None
    QDRANT_COLLECTION_NAME: str = _require_env("QDRANT_COLLECTION_NAME")
    QDRANT_USER_COLLECTION_NAME: str = _require_env("QDRANT_USER_COLLECTION_NAME")

    # RabbitMQ configuration (for background workers)
    RABBITMQ_URL: str = _require_env("RABBITMQ_URL")
    PRODUCT_CREATED_QUEUE: str = _require_env("PRODUCT_CREATED_QUEUE")
    USER_VECTOR_QUEUE: str = _require_env("USER_VECTOR_QUEUE")

    # Embedding model configuration
    EMBEDDING_MODEL_NAME: str = _require_env("EMBEDDING_MODEL_NAME")


@lru_cache
def get_settings() -> Settings:
    return Settings()

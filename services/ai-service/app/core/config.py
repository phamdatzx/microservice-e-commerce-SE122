import json
import os
from pathlib import Path
from functools import lru_cache
from typing import Any

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


def _parse_positive_int_env(name: str, default: int) -> int:
    raw = os.getenv(name)
    if raw is None:
        return default
    value = int(raw.strip())
    if value <= 0:
        raise RuntimeError(f"{name} must be a positive integer, got {value!r}")
    return value


def _parse_interaction_action_scores_json() -> dict[str, float]:
    """
    JSON object mapping interaction action -> weight, e.g.
    {"view": 1, "click": 2, "purchase": 10}
    """
    raw = os.getenv("INTERACTION_ACTION_SCORES_JSON", '{"view": 1}')
    try:
        data: Any = json.loads(raw)
    except json.JSONDecodeError as exc:
        raise RuntimeError(
            f"Invalid INTERACTION_ACTION_SCORES_JSON (must be JSON object): {exc}"
        ) from exc
    if not isinstance(data, dict):
        raise RuntimeError("INTERACTION_ACTION_SCORES_JSON must be a JSON object")
    return {str(k): float(v) for k, v in data.items()}


def _parse_float_env(name: str, default: float) -> float:
    raw = os.getenv(name)
    if raw is None:
        return default
    return float(raw.strip())


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

    # Intent model configuration
    INTENT_MODEL_REPO_ID: str = _require_env("INTENT_MODEL_REPO_ID")
    INTENT_DATASET_URL: str = _require_env("INTENT_DATASET_URL")

    # MongoDB (user interaction history + workers that persist to Mongo)
    MONGO_URI: str | None = os.getenv("MONGO_URI") or None
    MONGO_DB_NAME: str | None = os.getenv("MONGO_DB_NAME") or None
    INTERACTION_COLLECTION_NAME: str = os.getenv(
        "INTERACTION_COLLECTION_NAME", "user_interactions"
    )

    # Newest N interactions used to compute user vector (per user)
    USER_INTERACTION_HISTORY_LIMIT: int = _parse_positive_int_env(
        "USER_INTERACTION_HISTORY_LIMIT", 100
    )

    # Score per action when message has no explicit `score`
    INTERACTION_ACTION_SCORES: dict[str, float] = _parse_interaction_action_scores_json()
    INTERACTION_DEFAULT_SCORE: float = _parse_float_env("INTERACTION_DEFAULT_SCORE", 1.0)


@lru_cache
def get_settings() -> Settings:
    return Settings()

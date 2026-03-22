from functools import lru_cache
from typing import List

from sentence_transformers import SentenceTransformer

from app.core.config import get_settings
from app.schemas.product_vector import ProductPayload


@lru_cache
def get_embedding_model() -> SentenceTransformer:
    settings = get_settings()
    return SentenceTransformer(settings.EMBEDDING_MODEL_NAME)


def _build_product_text(payload: ProductPayload) -> str:
    """
    Build a textual representation of the product for embedding.
    You can tune this later (add description, attributes, etc.).
    """
    parts: List[str] = [
        payload.name,
        payload.category_name,
    ]

    return " | ".join(parts)


def compute_product_embedding(payload: ProductPayload) -> List[float]:
    model = get_embedding_model()
    text = _build_product_text(payload)
    vector = model.encode(text, normalize_embeddings=True)
    return vector.tolist()


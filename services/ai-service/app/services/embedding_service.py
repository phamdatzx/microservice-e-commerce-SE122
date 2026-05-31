"""
Product embedding service.

Uses OpenAI ``text-embedding-3-small`` (1536-dim) to embed product text,
the same model used by the RAG pipeline, so both the ``products`` and
``product_docs`` Qdrant collections share a single embedding space.
"""

from __future__ import annotations

from functools import lru_cache
from typing import List

from openai import OpenAI

from app.core.config import get_settings
from app.schemas.product_vector import ProductPayload


@lru_cache
def _get_openai_client() -> OpenAI:
    settings = get_settings()
    return OpenAI(api_key=settings.OPENAI_API_KEY)


def _build_product_text(payload: ProductPayload) -> str:
    """
    Build a textual representation of the product for embedding.
    Mirrors the text format used in ``scripts/embed_products.py``.
    """
    parts: List[str] = [payload.name, payload.category_name]
    return " | ".join(p for p in parts if p)


def compute_product_embedding(payload: ProductPayload) -> List[float]:
    """Return a 1536-dim unit-norm embedding for the given product payload."""
    settings = get_settings()
    client = _get_openai_client()
    text = _build_product_text(payload)

    response = client.embeddings.create(
        model=settings.OPENAI_EMBEDDING_MODEL,
        input=text,
    )
    return response.data[0].embedding

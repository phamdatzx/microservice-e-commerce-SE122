"""
RAG retriever service — connects to the Qdrant ``product_docs`` collection
(populated by ``scripts/embed_products.py``) and returns a LangChain-
compatible retriever that embeds queries with **text-embedding-3-small**
(same model used at indexing time).

Public API
----------
- ``get_retriever(k=...)``                    → LangChain ``VectorStoreRetriever``
- ``retrieve_products(query, k=...)``         → list of ``Document`` (plain similarity)
- ``retrieve_products_filtered(query, ...)``  → list of ``Document`` (with Qdrant payload filters)
"""

from __future__ import annotations

import logging
from functools import lru_cache
from typing import Any

from langchain_openai import OpenAIEmbeddings
from langchain_qdrant import QdrantVectorStore
from langchain_core.vectorstores import VectorStoreRetriever
from langchain_core.documents import Document
from qdrant_client import QdrantClient
from qdrant_client.http import models as qmodels

from app.core.config import get_settings

logger = logging.getLogger(__name__)


# ---------------------------------------------------------------------------
# Singleton clients — created once per process
# ---------------------------------------------------------------------------

@lru_cache
def _get_qdrant_client() -> QdrantClient:
    settings = get_settings()
    return QdrantClient(url=settings.QDRANT_URL, api_key=settings.QDRANT_API_KEY)


@lru_cache
def _get_embeddings() -> OpenAIEmbeddings:
    """Return an OpenAI embedding instance using **text-embedding-3-small**
    (must match the model used in ``scripts/embed_products.py``)."""
    settings = get_settings()
    return OpenAIEmbeddings(
        model="text-embedding-3-small",
        openai_api_key=settings.OPENAI_API_KEY,
    )


@lru_cache
def _get_vector_store() -> QdrantVectorStore:
    """Build a LangChain QdrantVectorStore backed by the existing
    ``product_docs`` collection."""
    settings = get_settings()
    client = _get_qdrant_client()
    embeddings = _get_embeddings()

    vector_store = QdrantVectorStore(
        client=client,
        collection_name=settings.QDRANT_PRODUCT_DOC_COLLECTION,
        embedding=embeddings,
    )

    logger.info(
        "Connected to Qdrant collection '%s' for RAG retrieval.",
        settings.QDRANT_PRODUCT_DOC_COLLECTION,
    )
    return vector_store


# ---------------------------------------------------------------------------
# Filter builder
# ---------------------------------------------------------------------------

def build_qdrant_filter(
    *,
    price_max: int | None = None,
    price_min: int | None = None,
    in_stock: bool = True,
    is_active: bool = True,
    min_rating: float | None = None,
) -> qmodels.Filter | None:
    """Build a Qdrant ``Filter`` with ``must`` conditions.

    Parameters
    ----------
    price_max : int, optional
        Maximum price (filters ``price_min <= price_max``).
    price_min : int, optional
        Minimum price (filters ``price_max >= price_min``).
    in_stock : bool
        If True, only return products with ``stock > 0``.
    is_active : bool
        If True, only return products with ``is_active == true``.
    min_rating : float, optional
        Minimum average rating (e.g. 4.0).

    Returns
    -------
    qmodels.Filter or None
        A Qdrant filter object, or None if no conditions apply.
    """
    must: list[qmodels.Condition] = []

    if is_active:
        must.append(
            qmodels.FieldCondition(
                key="is_active",
                match=qmodels.MatchValue(value=True),
            )
        )

    if in_stock:
        must.append(
            qmodels.FieldCondition(
                key="stock",
                range=qmodels.Range(gt=0),
            )
        )

    if price_max is not None:
        # Product's minimum price should be ≤ the user's budget
        must.append(
            qmodels.FieldCondition(
                key="price_min",
                range=qmodels.Range(lte=price_max),
            )
        )

    if price_min is not None:
        # Product's maximum price should be ≥ the user's floor
        must.append(
            qmodels.FieldCondition(
                key="price_max",
                range=qmodels.Range(gte=price_min),
            )
        )

    if min_rating is not None:
        must.append(
            qmodels.FieldCondition(
                key="rating",
                range=qmodels.Range(gte=min_rating),
            )
        )

    if not must:
        return None

    return qmodels.Filter(must=must)


# ---------------------------------------------------------------------------
# Public API
# ---------------------------------------------------------------------------

def get_retriever(k: int | None = None) -> VectorStoreRetriever:
    """Return a LangChain retriever over the product-docs collection.

    Parameters
    ----------
    k : int, optional
        Number of documents to retrieve.  Defaults to the ``RAG_TOP_K``
        setting (typically 5).

    Returns
    -------
    VectorStoreRetriever
        Ready-to-use retriever that can be plugged into any LangChain chain.
    """
    settings = get_settings()
    top_k = k if k is not None else settings.RAG_TOP_K

    vector_store = _get_vector_store()
    retriever = vector_store.as_retriever(
        search_type="similarity",
        search_kwargs={"k": top_k},
    )

    logger.debug("Created retriever with k=%d", top_k)
    return retriever


def retrieve_products(query: str, k: int | None = None) -> list[Document]:
    """Convenience wrapper: embed *query* and return the top-k product
    documents from Qdrant (no payload filtering).

    Parameters
    ----------
    query : str
        Natural-language user query (e.g. "iPhone giá rẻ dưới 10 triệu").
    k : int, optional
        Number of documents.  Defaults to ``RAG_TOP_K``.

    Returns
    -------
    list[Document]
        Each document has ``.page_content`` (the text block) and
        ``.metadata`` (full product payload stored at indexing time).
    """
    retriever = get_retriever(k=k)
    docs = retriever.invoke(query)
    logger.info(
        "Retrieved %d documents for query: '%.80s…'",
        len(docs),
        query,
    )
    return docs


def retrieve_products_filtered(
    query: str,
    *,
    k: int | None = None,
    price_max: int | None = None,
    price_min: int | None = None,
    in_stock: bool = True,
    is_active: bool = True,
    min_rating: float | None = None,
) -> list[Document]:
    """Embed *query* and return top-k products with Qdrant payload filtering.

    This uses the raw Qdrant client (not the LangChain retriever) to pass
    structured filters alongside the vector search.

    Parameters
    ----------
    query : str
        Semantic search query — should be product keywords only
        (e.g. "iPhone", "máy giặt"), NOT price or other filters.
    k : int, optional
        Number of results.  Defaults to ``RAG_TOP_K``.
    price_max : int, optional
        Upper price bound (VND).
    price_min : int, optional
        Lower price bound (VND).
    in_stock : bool
        Only return in-stock products (default True).
    is_active : bool
        Only return active products (default True).
    min_rating : float, optional
        Minimum average star rating.

    Returns
    -------
    list[Document]
        Documents with ``.page_content`` and ``.metadata``.
    """
    settings = get_settings()
    top_k = k if k is not None else settings.RAG_TOP_K

    # Embed the query
    embeddings = _get_embeddings()
    query_vector = embeddings.embed_query(query)

    # Build filter
    qdrant_filter = build_qdrant_filter(
        price_max=price_max,
        price_min=price_min,
        in_stock=in_stock,
        is_active=is_active,
        min_rating=min_rating,
    )

    # Search Qdrant directly
    client = _get_qdrant_client()
    results = client.query_points(
        collection_name=settings.QDRANT_PRODUCT_DOC_COLLECTION,
        query=query_vector,
        query_filter=qdrant_filter,
        limit=top_k,
        with_payload=True,
    )

    # Convert to LangChain Documents
    points = results.points if hasattr(results, "points") else results
    docs: list[Document] = []
    for point in points:
        payload: dict[str, Any] = point.payload or {}
        # Reconstruct page_content from the semantic fields
        content_parts = [payload.get("name", "")]
        desc = payload.get("description", "")
        if desc:
            content_parts.append(desc)
        cats = payload.get("category_names", [])
        if cats:
            content_parts.append(f"Category: {', '.join(cats)}")

        docs.append(
            Document(
                page_content="\n".join(content_parts),
                metadata=payload,
            )
        )

    logger.info(
        "Retrieved %d filtered documents for query='%.60s' "
        "(price_min=%s, price_max=%s, in_stock=%s, min_rating=%s)",
        len(docs), query, price_min, price_max, in_stock, min_rating,
    )
    return docs

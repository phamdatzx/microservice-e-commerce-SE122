"""
RAG retriever service — connects to the Qdrant ``product_docs`` collection
(populated by ``scripts/embed_products.py``) and returns a LangChain-
compatible retriever that embeds queries with **text-embedding-ada-002**
(same model used at indexing time).

Public API
----------
- ``get_retriever(k=...)`` → LangChain ``VectorStoreRetriever``
- ``retrieve_products(query, k=...)`` → list of ``Document`` objects
"""

from __future__ import annotations

import logging
from functools import lru_cache

from langchain_openai import OpenAIEmbeddings
from langchain_qdrant import QdrantVectorStore
from langchain_core.vectorstores import VectorStoreRetriever
from langchain_core.documents import Document
from qdrant_client import QdrantClient

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
    """Return an OpenAI embedding instance using **text-embedding-ada-002**
    (must match the model used in ``scripts/embed_products.py``)."""
    settings = get_settings()
    return OpenAIEmbeddings(
        model="text-embedding-ada-002",
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
    documents from Qdrant.

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

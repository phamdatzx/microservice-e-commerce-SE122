from functools import lru_cache

from qdrant_client import QdrantClient
from qdrant_client.http import models as qmodels

from app.core.config import get_settings
from app.schemas.product_vector import ProductVector


@lru_cache
def get_qdrant_client() -> QdrantClient:
    settings = get_settings()
    return QdrantClient(url=settings.QDRANT_URL, api_key=settings.QDRANT_API_KEY)


def _ensure_collection(client: QdrantClient, collection_name: str, vector_size: int) -> None:
    collections = client.get_collections().collections or []
    if any(c.name == collection_name for c in collections):
        return

    client.create_collection(
        collection_name=collection_name,
        vectors_config=qmodels.VectorParams(
            size=vector_size,
            distance=qmodels.Distance.COSINE,
        ),
    )


def upsert_product_vector(product: ProductVector) -> None:
    settings = get_settings()
    client = get_qdrant_client()

    vector_size = len(product.vector)
    if vector_size == 0:
        raise ValueError("Vector must not be empty")

    _ensure_collection(client, settings.QDRANT_COLLECTION_NAME, vector_size)

    client.upsert(
        collection_name=settings.QDRANT_COLLECTION_NAME,
        points=[
            qmodels.PointStruct(
                id=product.id,
                vector=product.vector,
                payload=product.payload.model_dump(),
            )
        ],
    )


import os
from qdrant_client import QdrantClient
from qdrant_client.models import Distance, VectorParams

QDRANT_HOST = os.getenv("QDRANT_HOST", "localhost")
QDRANT_PORT = int(os.getenv("QDRANT_PORT", 6333))
QDRANT_PRODUCT_COLLECTION = os.getenv("QDRANT_PRODUCT_COLLECTION", "products")
QDRANT_RATING_COLLECTION = os.getenv("QDRANT_RATING_COLLECTION", "ratings")

qc = QdrantClient(host=QDRANT_HOST, port=QDRANT_PORT)

print("Đang xoá và tạo lại collection với kích thước Vector là 1024 (BAAI/bge-m3)...")

qc.recreate_collection(
    collection_name=QDRANT_PRODUCT_COLLECTION,
    vectors_config=VectorParams(size=1024, distance=Distance.COSINE)
)
print("-> Đã tạo lại collection:", QDRANT_PRODUCT_COLLECTION)

qc.recreate_collection(
    collection_name=QDRANT_RATING_COLLECTION,
    vectors_config=VectorParams(size=1024, distance=Distance.COSINE)
)
print("-> Đã tạo lại collection:", QDRANT_RATING_COLLECTION)

print("HOÀN TẤT!")

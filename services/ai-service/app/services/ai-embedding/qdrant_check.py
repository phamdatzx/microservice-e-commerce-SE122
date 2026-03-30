from qdrant_client import QdrantClient
qc = QdrantClient(host="localhost", port=6333)
print(qc.get_collection("docs"))
print(qc.scroll(collection_name="docs", limit=10))
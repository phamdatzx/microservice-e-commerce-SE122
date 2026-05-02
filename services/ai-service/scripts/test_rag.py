"""Quick smoke test for rag_service.retrieve_products()."""

import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).resolve().parents[1]))

from app.services.rag_service import retrieve_products

if __name__ == "__main__":
    query = "toôi muốn mua máy giặt"
    k = 5

    print(f"Query: {query}")
    print(f"Top-k: {k}")
    print("=" * 60)

    docs = retrieve_products(query, k=k)

    if not docs:
        print("⚠️  No documents returned.")
    else:
        for i, doc in enumerate(docs, 1):
            print(f"\n--- Result {i} ---")
            print(f"Content:\n{doc.page_content[:300]}")
            print(f"\nMetadata: {doc.metadata}")
            print()

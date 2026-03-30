1. Tạo virtualenv và cài dependencies:
   python -m venv .venv
   .venv\Scripts\activate
   pip install -r requirements.txt

2. Chạy RabbitMQ và Qdrant (Docker Desktop recommended):
   docker run -d --name rabbitmq -p 5672:5672 -p 15672:15672 rabbitmq:3.11-management
   docker run -d --name qdrant -p 6333:6333 qdrant/qdrant:1.16.3

3. Copy .env.example -> .env và chỉnh nếu cần.

4. Chạy consumer:
   python consumer_sync.py

5. Test gửi messages:
   python producer_test.py

6. Kiểm tra Qdrant:
   from qdrant_client import QdrantClient
   qc = QdrantClient(host="localhost", port=6333)
   print(qc.get_collection("docs"))
   print(qc.scroll(collection_name="docs", limit=10))

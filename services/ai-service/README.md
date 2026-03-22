# AI Service (FastAPI)

Recommended FastAPI structure for scalable development.

## Folder structure

```text
app/
  api/
    v1/
      endpoints/
      router.py
  core/
    application.py
    config.py
  models/
  schemas/
  services/
  main.py
```

## Run locally

```bash
uvicorn app.main:app --reload
```

## Qdrant configuration

Environment variables used to connect to Qdrant:

- `QDRANT_URL` – Qdrant HTTP URL (default `http://localhost:6333`)
- `QDRANT_API_KEY` – Qdrant API key (optional, for Qdrant Cloud)
- `QDRANT_COLLECTION_NAME` – collection name for product vectors (default `products`)

Example:

```bash
export QDRANT_URL="http://localhost:6333"
export QDRANT_COLLECTION_NAME="products"
uvicorn app.main:app --reload
```

## User vector worker + MongoDB interaction history

`app/user_vector_worker.py` consumes `USER_VECTOR_QUEUE`, appends each event to MongoDB, then recomputes the user’s vector in Qdrant from the **newest N** interactions (weighted by `score`).

Required for the worker:

- `MONGO_URI` – MongoDB connection string
- `MONGO_DB_NAME` – database name (e.g. same as `product-service` `DB_NAME`)

Optional:

- `INTERACTION_COLLECTION_NAME` – collection for documents (default `user_interactions`)
- `USER_INTERACTION_HISTORY_LIMIT` – how many latest interactions to use (default `100`)
- `INTERACTION_ACTION_SCORES_JSON` – JSON object mapping `action` → weight when the message omits `score`, e.g. `{"view":1,"purchase":10}`
- `INTERACTION_DEFAULT_SCORE` – weight for unknown actions (default `1`)

RabbitMQ message (JSON):

```json
{
  "user_id": "u123",
  "product_id": "p456",
  "action": "view",
  "score": 1,
  "timestamp": "2026-03-22T20:00:00Z"
}
```

`score` and `timestamp` are optional; if `score` is omitted, it is resolved from `action` using `INTERACTION_ACTION_SCORES_JSON`.

Run the worker (from `services/ai-service`):

```bash
python -m app.user_vector_worker
```

## Next steps

- Add feature modules under `app/api/v1/endpoints/` (for example: `chat.py`, `embeddings.py`)
- Put business logic in `app/services/`
- Put request/response models in `app/schemas/`
- Add persistence integration in `app/models/` when needed

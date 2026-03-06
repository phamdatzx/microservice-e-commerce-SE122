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

## Next steps

- Add feature modules under `app/api/v1/endpoints/` (for example: `chat.py`, `embeddings.py`)
- Put business logic in `app/services/`
- Put request/response models in `app/schemas/`
- Add persistence integration in `app/models/` when needed

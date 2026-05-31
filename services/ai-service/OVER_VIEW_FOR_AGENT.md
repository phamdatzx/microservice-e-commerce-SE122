# AI Service — Codebase Overview

## Tech Stack

| Layer | Technology |
|---|---|
| Framework | FastAPI (Python) |
| LLM | GPT-4o via OpenAI API |
| Orchestration | LangChain (`langchain-classic`, `langchain-openai`, `langchain-qdrant`) |
| Vector DB | Qdrant (local Docker or Qdrant Cloud) |
| Message Queue | RabbitMQ (`pika`) |
| Document DB | MongoDB Atlas (`pymongo`) — interaction history + CF results |
| Embedding (products) | `intfloat/multilingual-e5-small` (sentence-transformers) |
| Embedding (RAG docs) | `text-embedding-3-small` via OpenAI |
| Intent Classifier | PhoBERT-based model (`duc004505/e-commerce-intent-classifier-model`) from HuggingFace |
| NER model | `duc004505/ner_model_output` from HuggingFace |
| Vietnamese NLP | `underthesea` (word tokenization) |
| Numerical | numpy, scipy (sparse), scikit-learn (cosine similarity), pandas |

---

## Architecture Overview

```
                                ┌──────────────────────────────────────────────────┐
                                │                 ai-service (FastAPI)             │
                                │                                                  │
 Client ──── POST /chat ───────▶│  intent_router → predict_intent → (agent/LLM)   │
                                │                                                  │
 Client ── POST /users/recs ───▶│  content-based → Qdrant user_vectors query       │
                                │                                                  │
 Client ── GET /cf/{pid} ──────▶│  CF lookup → MongoDB item_similarity              │
                                │                                                  │
 Client ── POST /intent/pred ──▶│  predict_intent → PhoBERT classifier             │
                                └──────────────────────────────────────────────────┘
         RabbitMQ queues:
           product.created ──▶ product_embed_worker → Qdrant (products collection)
           user.interaction ──▶ user_vector_worker → MongoDB + Qdrant (user_vectors)
         Batch job (manual/cron):
           cf_worker ──▶ MongoDB interactions → item-item similarity → MongoDB item_similarity
```

---

## API Endpoints (`/api/ai`)

| Method | Path | Description |
|--------|------|-------------|
| `GET` | `/health` | Service health check |
| `POST` | `/chat` | Main chatbot entry point — routes via intent, then calls agent |
| `POST` | `/users/recommendations` | Content-based product recommendations for a user |
| `GET` | `/recommend/cf/{product_id}` | Item-based CF recommendations for a product |
| `POST` | `/intent/predict` | Raw intent classification (standalone endpoint) |
| `POST` | `/users/vector` | Upsert/compute user vector manually |

---

## Module Breakdown

### `app/agent/agent.py` — The LangChain Agent (Core Chatbot)

- **Architecture:** `AgentExecutor` with `create_tool_calling_agent` (OpenAI function calling)
- **Model:** GPT-4o (`temperature=0`)
- **System prompt:** Vietnamese shopping assistant
- **Tools available to the LLM:**
  1. `search_products` — RAG via Qdrant `product_docs` collection (semantic search + payload filters)
  2. `get_product_by_id` — (MOCK) product detail lookup
  3. `get_product_reviews` — (MOCK) product reviews
  4. `get_my_orders` — (MOCK) user orders
  5. `get_seller_info_by_id` — (MOCK) seller info
  6. `get_my_vouchers` — (MOCK) user vouchers
- **Important:** Tools 2–6 in `app/tools/backend_tools.py` are **currently mocked** — need real HTTP calls to other microservices
- `search_products` is **real and working** (hits Qdrant)

### `app/services/rag_service.py` — RAG Retrieval

- Uses **OpenAI `text-embedding-3-small`** to embed queries (must match indexing model)
- Queries Qdrant `product_docs` collection
- Supports payload filters: `price_min`, `price_max`, `in_stock`, `is_active`, `min_rating`
- Two access paths:
  - `get_retriever()` — LangChain VectorStoreRetriever (simple similarity)
  - `retrieve_products_filtered()` — Raw Qdrant client with structured filters (used by agent)

### `app/services/recommendation_service.py` — Content-Based Recommendations

- Fetches user vector from Qdrant `user_vectors` collection
- Queries Qdrant `products` collection for nearest product vectors
- Fast, real-time, O(log N) via HNSW index

### `app/services/cf_service.py` — Collaborative Filtering

- **Item-based CF** using cosine similarity
- Pipeline: load 30-day interactions from MongoDB → sparse user-item matrix → row-by-row cosine similarity → store top-K per product in MongoDB `item_similarity`
- Memory-efficient: never materializes full N×N matrix
- Filters: min user/item interactions, max 5000 products cap

### `app/services/user_vector_service.py` — User Vector Computation

- Fetches product embeddings from Qdrant `products` collection
- Computes weighted average of product vectors (weights = interaction scores)
- Normalizes to unit length
- Upserts result to Qdrant `user_vectors` collection

### `app/services/intent_router.py` — Intent Routing

- Calls `predict_intent()` to classify user query
- **Currently incomplete** — routing logic is commented out, returns intent but doesn't dispatch to different handlers yet

### `app/services/predict_intent_service.py` — Intent Classification (API)

- Simple wrapper that calls the PhoBERT intent model

### `app/services/predict_intent_and_keyword_service.py` — Full NLP Pipeline

- Loads both **Intent Classifier** (PhoBERT) + **NER model** from HuggingFace
- Segments Vietnamese text with `underthesea.word_tokenize`
- Outputs: intent label + confidence + extracted entities (product name, category, etc.)
- Intent classes: `search_product`, `product_detail`, `review_summary`, `compare_product`, `check_order_status`, `product_voucher`, `seller_voucher`, `product_applicable_voucher`, `seller_category`, `greeting`, `goodbye`, `out_of_scope`

---

## Workers

| Worker | Queue/Trigger | Function |
|--------|---------------|----------|
| `product_embed_worker.py` | `product.created` (RabbitMQ) | Computes product embedding (`intfloat/multilingual-e5-small`) → upserts to Qdrant `products` |
| `user_vector_worker.py` | `user.interaction` (RabbitMQ) | Persists interaction to MongoDB → recomputes user vector → upserts to Qdrant `user_vectors` |
| `cf_worker.py` | Manual / cron (batch, one-shot) | Full CF pipeline: MongoDB → sparse matrix → cosine similarity → MongoDB `item_similarity` |

---

## Qdrant Collections

| Collection | Embedding Model | Content |
|---|---|---|
| `products` | `intfloat/multilingual-e5-small` (384-dim) | Product semantic vectors (name + category) |
| `product_docs` | `text-embedding-3-small` (OpenAI) | RAG document chunks (name + description + category) |
| `user_vectors` | Same dim as `products` | Per-user preference vectors (weighted avg of interacted products) |

> **Critical mismatch note:** The `products` and `product_docs` collections use **different embedding models**. Product vectors for CB filtering use sentence-transformers; RAG uses OpenAI embeddings. Do not mix these.

---

## MongoDB Collections

| Collection | Purpose |
|---|---|
| `user_interactions` | Raw interaction history (`user_id`, `product_id`, `action`, `score`, `timestamp`) |
| `item_similarity` | Pre-computed CF results: `{product_id, similar_items: [{product_id, score}], updated_at}` |

---

## Scripts (`/scripts`)

| Script | Purpose |
|---|---|
| `embed_products.py` | Bulk-embed existing products into Qdrant `product_docs` (first-time setup) |
| `product_backfill.py` | Backfill products into Qdrant `products` collection |
| `reset_qdrant.py` | Delete and recreate Qdrant collections |
| `send_message.py` | Test utility to send RabbitMQ messages |
| `test_agent.py` | Interactive CLI test for the LangChain agent |
| `test_rag.py` | Test RAG retrieval |

---

## Current Implementation Status

| Feature | Status |
|---|---|
| Content-based recommendation | ✅ Complete |
| Item-based CF recommendation | ✅ Complete |
| Product embedding worker | ✅ Complete |
| User vector worker | ✅ Complete |
| RAG search (filtered) | ✅ Complete |
| LangChain agent (GPT-4o) | ✅ Complete |
| Intent classification (PhoBERT) | ✅ Complete (model loads from HuggingFace) |
| NER extraction | ✅ Complete (model loads from HuggingFace) |
| Intent routing in `/chat` | ⚠️ Incomplete — routing logic commented out, always falls through |
| Backend tools (orders, vouchers, etc.) | ⚠️ All **MOCKED** — need real HTTP calls to other services |
| Multi-turn chat history | ⚠️ Agent supports it but `/chat` endpoint is single-turn only |

---

## Key Design Decisions

1. **Two separate embedding spaces** — product embeddings (sentence-transformers, fast, local) vs RAG embeddings (OpenAI, richer). CB filtering uses the former; chatbot RAG uses the latter.
2. **CF is fully pre-computed** — no computation at request time; fast O(1) MongoDB lookup.
3. **User vector = weighted normalized average** of product embeddings. Interaction scores weight: view=1, click=2, add_to_cart=10, purchase=10.
4. **Agent tools are mocked** — backend_tools.py is a clear integration point. Each tool needs to be replaced with an HTTP call to the corresponding microservice.
5. **`intent_router.py` is the main integration gap** — it classifies intent but doesn't yet dispatch to the right handler (agent vs. direct API vs. CF).

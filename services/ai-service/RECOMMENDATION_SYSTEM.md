# Recommendation System — Architecture & Algorithm Documentation

## 1. System Architecture Overview

The recommendation system lives inside the **ai-service** (FastAPI / Python) and produces two independent recommendation signals that can be combined into a hybrid approach.

```
┌──────────────┐    events     ┌─────────────┐    RabbitMQ     ┌──────────────────────┐
│ product-     │──────────────▶│  RabbitMQ    │───────────────▶│  Workers (Python)    │
│ service (Go) │               │  Message     │                │                      │
│              │               │  Queue       │                │  ┌─────────────────┐ │
│ order-       │──────────────▶│              │                │  │ worker.py        │ │
│ service (Go) │               └──────────────┘                │  │ (product embed)  │ │
└──────────────┘                                               │  ├─────────────────┤ │
                                                               │  │ user_vector_     │ │
                                                               │  │ worker.py        │ │
                                                               │  │ (interaction →   │ │
                                                               │  │  user vector)    │ │
                                                               │  └─────────────────┘ │
                                                               └──────────┬───────────┘
                                                                          │ writes
                                                          ┌───────────────┼───────────────┐
                                                          ▼               ▼               ▼
                                                     ┌─────────┐   ┌──────────┐   ┌────────────┐
                                                     │ Qdrant   │   │ MongoDB  │   │ MongoDB    │
                                                     │ (vectors)│   │ user_    │   │ item_      │
                                                     │          │   │ interact │   │ similarity │
                                                     └────┬─────┘   └──────────┘   └─────┬──────┘
                                                          │                              │
                                                          ▼                              ▼
                                                   ┌─────────────────────────────────────────┐
                                                   │          FastAPI  (ai-service)           │
                                                   │                                         │
                                                   │  POST /users/recommendations             │
                                                   │       → Content-based (Qdrant)           │
                                                   │                                         │
                                                   │  GET  /recommend/cf/{product_id}         │
                                                   │       → Collaborative Filtering (Mongo)  │
                                                   └─────────────────────────────────────────┘

                              ┌─────────────────────────────────────────────┐
                              │         cf_worker.py  (batch / cron)        │
                              │  Reads MongoDB interactions → computes      │
                              │  item-item similarity → writes to MongoDB   │
                              └─────────────────────────────────────────────┘
```

### Component Responsibilities

| Component | Role |
|-----------|------|
| **product-service / order-service** (Go) | Publish user interaction events (`view`, `click`, `add_to_cart`, `purchase`) to RabbitMQ |
| **worker.py** | Consumes `product.created` queue → computes product embedding → stores vector in Qdrant |
| **user_vector_worker.py** | Consumes `user.interaction` queue → persists interaction in MongoDB → recomputes user vector → upserts into Qdrant |
| **cf_worker.py** | **Batch job** (cron / manual). Reads interaction history from MongoDB → computes item-item cosine similarity → stores top-K results in MongoDB `item_similarity` collection |
| **FastAPI endpoints** | Serve recommendations to clients |

---

## 2. Recommendation Approaches

### 2.1 Content-Based Filtering (existing)

**Idea:** Recommend products whose embeddings are close to the user's aggregated preference vector.

**Flow:**
1. Each product is embedded into a dense vector using a sentence-transformer model (`intfloat/multilingual-e5-small`). The embedding captures semantic meaning of the product name + category.
2. When a user interacts with products, a **user vector** is computed as a weighted average of the product vectors, where weights are interaction scores.
3. At recommendation time, the user vector is used to query Qdrant for the nearest product vectors (cosine similarity in vector space).

**Endpoint:** `POST /users/recommendations`

```
User Vector = Σ(score_i × product_vector_i) / Σ(score_i)
              normalized to unit length
```

### 2.2 Collaborative Filtering — Item-Based (new)

**Idea:** If many users who interacted with product A also interacted with product B, then A and B are similar. Recommend B when a user is looking at A.

**Flow:**
1. Load interaction history from MongoDB (last 30 days)
2. Build a user–item interaction matrix
3. Compute item-item cosine similarity row-by-row
4. Store top-K similar items per product in MongoDB
5. Serve pre-computed results via API

**Endpoint:** `GET /recommend/cf/{product_id}?limit=10`

---

## 3. Collaborative Filtering — Algorithm Detail

### 3.1 Data Loading & Filtering

**Source:** MongoDB `user_interactions` collection.

Each document:
```json
{
  "user_id": "u123",
  "product_id": "p456",
  "action": "purchase",
  "score": 10.0,
  "timestamp": "2026-03-20T10:00:00Z"
}
```

The `score` field is pre-resolved at insert time using the environment variable mapping:
```
INTERACTION_ACTION_SCORES_JSON = {"view": 1, "click": 2, "purchase": 10, "add_to_cart": 10}
```

**Filtering pipeline (in order):**

| Step | Filter | Default | Purpose |
|------|--------|---------|---------|
| 1 | Time window | Last 30 days | Keep only recent, relevant interactions |
| 2 | Min user interactions | `CF_MIN_USER_INTERACTIONS = 2` | Users with only 1 interaction provide no co-occurrence signal |
| 3 | Min item interactions | `CF_MIN_ITEM_INTERACTIONS = 2` | Items with too few interactions have unreliable similarity |
| 4 | Max products cap | `CF_MAX_PRODUCTS = 5000` | Keep only top-N products by interaction count to prevent OOM |

### 3.2 Matrix Construction

**Goal:** Build a sparse User × Item matrix where each cell = aggregated interaction score.

```
              Product_A  Product_B  Product_C  ...
  User_1         2.0        0.0       10.0
  User_2         1.0        3.0        0.0
  User_3         0.0       10.0       10.0
  ...
```

**Steps:**

1. **Aggregate** — If a user interacted with the same product multiple times (e.g., viewed then purchased), scores are summed:
   ```
   GROUP BY (user_id, product_id) → SUM(score)
   ```

2. **Encode** — Convert user_id and product_id strings to integer indices using pandas categorical encoding (0-indexed).

3. **Sparse matrix** — Construct a `scipy.sparse.csr_matrix` using (row_indices, col_indices, values). Only non-zero cells are stored in memory.

   ```python
   # Only stores nnz (non-zero) entries, not the full M×N grid
   sparse_mat = csr_matrix((data, (row, col)), shape=(n_users, n_items))
   ```

   **Memory comparison** for 5000 products × 10000 users:
   - Dense: 5000 × 10000 × 4 bytes = **200 MB**
   - Sparse (e.g. 50,000 non-zero entries): ~50,000 × 12 bytes ≈ **0.6 MB**

### 3.3 Similarity Computation

**Algorithm:** Cosine similarity, computed **row-by-row** (never materializing the full N×N matrix).

**Cosine similarity formula:**

```
                    A · B
sim(A, B) = ─────────────────
             ‖A‖ × ‖B‖
```

Where A and B are the interaction vectors for two items (each vector has one dimension per user).

**Why item-based (not user-based)?**
- The item × user matrix is transposed from the user × item matrix
- Each row represents an item's interaction pattern across all users
- Two items are similar if they are interacted with by the same set of users in similar proportions

**Row-by-row computation (memory safe):**

```
For each item i (row in item×user matrix):
    1. Take item_i's row vector (1 × n_users)  ← sparse, single row
    2. Compute cosine_similarity(item_i, ALL_items) → 1 × n_items result
    3. Set self-similarity to -1 (exclude self)
    4. Extract top-K highest scores using argpartition (O(N) vs O(N log N) for full sort)
    5. Keep only items with positive similarity
    6. Store result, discard the 1×N similarity vector
```

**Why this approach instead of full N×N matrix:**

| Approach | Memory | For N=5000 items |
|----------|--------|-------------------|
| Full N×N matrix | O(N²) | 5000² × 8 bytes = **200 MB** |
| Row-by-row | O(N) per iteration | 5000 × 8 bytes = **40 KB** at any time |

### 3.4 Top-K Selection

For each item, we keep only the `CF_TOP_K` (default: 20) most similar items.

**Efficient selection** using `numpy.argpartition`:
- `argpartition` runs in O(N) average time (vs O(N log N) for full sort)
- Only the selected top-K indices are then fully sorted (O(K log K))

Items with similarity ≤ 0 are discarded (no meaningful co-occurrence).

### 3.5 Storage (MongoDB)

Results are stored in the `item_similarity` collection:

```json
{
  "product_id": "p456",
  "similar_items": [
    { "product_id": "p789", "score": 0.923456 },
    { "product_id": "p012", "score": 0.871234 },
    ...
  ],
  "updated_at": "2026-03-29T14:00:00Z"
}
```

**Write strategy:**
- **Upsert** per product (bulk write in batches of 1000)
- **Delete stale** entries — products no longer in the current computation batch are removed
- A unique index on `product_id` ensures fast lookups and prevents duplicates

### 3.6 API Query

`GET /recommend/cf/{product_id}?limit=10`

Simply reads the pre-computed document from MongoDB and returns the top `limit` similar items. This is a fast O(1) lookup — no computation at request time.

---

## 4. Configuration Reference

All settings have defaults and can be overridden via environment variables or `.env`:

| Variable | Default | Description |
|----------|---------|-------------|
| `CF_TOP_K` | `20` | Number of similar items stored per product |
| `CF_MIN_USER_INTERACTIONS` | `2` | Skip users with fewer interactions |
| `CF_MIN_ITEM_INTERACTIONS` | `2` | Skip items with fewer interactions |
| `CF_MAX_PRODUCTS` | `5000` | Max products to include in matrix (by interaction count) |
| `ITEM_SIMILARITY_COLLECTION_NAME` | `item_similarity` | MongoDB collection name |
| `INTERACTION_ACTION_SCORES_JSON` | `{"view":1,"click":2,...}` | Score mapping per action type (shared with content-based) |

---

## 5. Running the CF Pipeline

```bash
# One-shot batch job (schedule via cron for periodic updates)
python -m app.workers.cf_worker
```

Example output:
```
2026-03-29 14:00:00 INFO app.services.cf_service: Loaded 15234 raw interactions from last 30 days
2026-03-29 14:00:00 INFO app.services.cf_service: After user filter (>=2 interactions): 12450 interactions
2026-03-29 14:00:00 INFO app.services.cf_service: After item filter (>=2 interactions): 11890 interactions
2026-03-29 14:00:00 INFO app.services.cf_service: After top-5000 products cap: 11890 interactions, 342 unique products
2026-03-29 14:00:00 INFO app.services.cf_service: Built sparse matrix: 1205 users × 342 items, nnz=11890
2026-03-29 14:00:01 INFO app.services.cf_service: Computed similarities for 340 items (with at least 1 similar item)
2026-03-29 14:00:01 INFO app.cf_worker: === CF batch job completed in 1.23 seconds. 340 products with similarities. ===
```

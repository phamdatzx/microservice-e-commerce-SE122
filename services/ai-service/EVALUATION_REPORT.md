# Recommendation System — Evaluation Report

This document explains the offline evaluation methodology, datasets, metrics, and results for the two recommendation algorithms implemented in this service.

---

## Table of Contents

1. [System Overview](#system-overview)
2. [Collaborative Filtering Evaluation](#1-collaborative-filtering-evaluation)
   - [Dataset](#dataset-movielens-100k)
   - [Methodology](#methodology)
   - [Metrics explained](#metrics-explained)
   - [Results](#results)
   - [How to run](#how-to-run-cf)
3. [Content-Based Filtering Evaluation](#2-content-based-filtering-evaluation)
   - [Dataset](#dataset-ecommerce_products_5kcsv)
   - [Methodology](#methodology-1)
   - [Metrics explained](#metrics-explained-1)
   - [Results](#results-1)
   - [How to run](#how-to-run-cb)
4. [Comparison & Conclusions](#comparison--conclusions)
5. [Limitations](#limitations)

---

## System Overview

| | Collaborative Filtering (CF) | Content-Based Filtering (CB) |
|---|---|---|
| **Algorithm** | Item-based CF, cosine similarity | OpenAI `text-embedding-3-small`, nearest-neighbor |
| **Signal** | User interaction history (implicit/explicit) | Product name + category text |
| **Storage** | MongoDB `item_similarity` | Qdrant `product_docs` |
| **API endpoint** | `GET /api/ai/recommend/cf/{product_id}` | `POST /api/ai/users/recommendations` |
| **Output** | "Products similar to X" (item→item) | "Products for user Y" (user→item) |
| **Update** | Batch job (`cf_worker`) | Real-time (embedding on product creation) |

---

## 1. Collaborative Filtering Evaluation

### Dataset: MovieLens 100k

| Property | Value |
|---|---|
| Source | [GroupLens Research, University of Minnesota](https://grouplens.org/datasets/movielens/100k/) |
| Ratings | 100,000 explicit ratings (1–5 stars) |
| Users | 943 |
| Movies (items) | 1,682 |
| Min ratings/user | 20 |
| Location in project | `ml-100k/` |

**Data mapping** — how MovieLens fields map to our system:

| MovieLens field | System field | Notes |
|---|---|---|
| `user_id` | `user_id` | direct mapping |
| `item_id` (movie) | `product_id` | movies treated as products |
| `rating` (1–5) | `score` | explicit ratings used as interaction weights |

This mapping is valid because our CF algorithm is purely numeric — it only uses the score values as weights in the cosine similarity computation. Whether a score represents a star rating or an implicit action score (view=1, purchase=10) does not change the algorithm.

### Methodology

The evaluation uses the **5-fold cross-validation** splits already provided by the MovieLens dataset (`u1.base`/`u1.test` through `u5.base`/`u5.test`). Each split is 80% training / 20% test.

**Pipeline per fold:**

```
u{n}.base (train)
    │
    ├─► build_user_item_matrix()       sparse matrix: users × movies
    │
    ├─► compute_item_similarity()      cosine similarity per item pair (top-20 neighbors)
    │
    └─► for each test user:
            training items  ──► fetch similar items for each
                            ──► aggregate scores, rank candidates
                            ──► exclude already-seen items
                            ──► top-N recommendation list
                                        │
u{n}.test (ground truth)               │
    items rated ≥ 4 = "relevant" ──────┘
                                        │
                                        ▼
                            Precision@K, Recall@K, NDCG@K, HitRate@K
```

**Average across all 5 folds** = final reported metrics.

**Relevance threshold:** a movie is considered "liked" if the user rated it **≥ 4 out of 5**.

**Random baseline:** with 1,682 movies and ~40 relevant per user in the test set, random Precision@10 ≈ **0.024** (2.4%).

### Metrics explained

| Metric | Formula | Meaning |
|---|---|---|
| **Precision@K** | hits@K / K | Of the K recommended movies, what fraction did the user actually like? |
| **Recall@K** | hits@K / \|relevant\| | Of all movies the user liked in the test set, what fraction appear in top-K? |
| **NDCG@K** | DCG@K / idealDCG@K | Are higher-rated movies ranked earlier? (1.0 = perfect ranking) |
| **HitRate@K** | 1 if hits@K ≥ 1 | Did at least one liked movie appear in the top-K list? |

> **NDCG (Normalized Discounted Cumulative Gain)** penalizes relevant items that appear lower in the list. A score of 0.33 at K=10 means the system is recovering about a third of the ideal ranking quality.

### Results

**5-Fold Cross-Validation** · `top_k=20` · `relevance ≥ 4` · ~743 avg evaluated users/fold

| K | Precision@K | Recall@K | NDCG@K | HitRate@K |
|--:|--:|--:|--:|--:|
| 5 | 0.2915 | 0.1459 | 0.3304 | 0.7094 |
| **10** | **0.2491** | **0.2284** | **0.3252** | **0.8134** |
| 20 | 0.2013 | 0.3396 | 0.3389 | 0.8943 |

**Key takeaways:**

- **HitRate@10 = 81%** — 4 out of 5 users receive at least one movie they actually liked in their top-10 list. This is the most practical metric for a recommendation system.
- **Precision@10 = 0.25** — one quarter of all recommended movies are relevant, which is **10× better than random** (2.4% baseline).
- **NDCG@10 = 0.33** — the system has a meaningful ability to rank relevant items higher, though there is room for improvement.
- As K grows, **Recall increases** (more coverage) while **Precision decreases** (more noise) — the expected trade-off.

### How to run (CF)

```bash
cd /path/to/ai-service

# Using the project venv
.venv/bin/python -m app.evaluation.cf_evaluator

# With custom options
.venv/bin/python -m app.evaluation.cf_evaluator \
  --data-dir ml-100k \
  --top-k 20 \
  --k-values 5,10,20 \
  --relevance-threshold 4

# Interactive notebook (requires jupyter)
.venv/bin/jupyter notebook cf_evaluation.ipynb
```

**Source:** `app/evaluation/cf_evaluator.py`

---

## 2. Content-Based Filtering Evaluation

### Dataset: ecommerce_products_5k.csv

| Property | Value |
|---|---|
| Source | Synthetic e-commerce dataset (generated for evaluation purposes) |
| Products | 5,000 |
| Sub-categories | 20 (250 products each, perfectly balanced) |
| Main categories | 10 (Electronics, Fashion, Home & Living, …) |
| Fields used | `product_id`, `name`, `sub_category` |
| Location in project | `data-contentbased/ecommerce_products_5k.csv` |
| Schema reference | `data-contentbased/structure.md` |

**Example products:**

| product_id | name | sub_category |
|---|---|---|
| PRD00001 | Apple Smartphones | Smartphones |
| PRD00251 | L'Oréal Skincare | Skincare |
| PRD01001 | Nike Footwear | Footwear |

### Methodology

The evaluation is **fully offline** — no Qdrant instance or MongoDB required. All similarity is computed directly in NumPy.

**Embed text format** (exact same format as the production system in `embedding_service.py`):

```
"{name} | {sub_category}"
```

Examples:
- `"Apple Smartphones | Smartphones"`
- `"Nike Footwear | Footwear"`

**Pipeline:**

```
ecommerce_products_5k.csv
    │
    ├─► build embed text: "{name} | {sub_category}"   (5,000 texts)
    │
    ├─► OpenAI text-embedding-3-small                 (1536-dim, same model as production)
    │       │
    │       └─► cache to .cache/cb_eval/embeddings.npy  (API called only once)
    │
    ├─► L2-normalise all vectors
    │
    └─► for each product (as query):
            dot product vs all 4,999 others (= cosine similarity on normalised vectors)
            take top-K (exclude self)
                    │
                    ▼
        ground truth: same sub_category = relevant (249 relevant items per query)
                    │
                    ▼
        Precision@K, Recall@K, HitRate@K, Mean Cosine Similarity
```

**Cost:** 5,000 × ~15 tokens ≈ 75,000 tokens → `text-embedding-3-small` at $0.02/1M tokens ≈ **$0.0015 total** (run once, cached forever).

**Random baseline:** with 250 relevant products per category out of 4,999 total, random Precision@10 ≈ **0.050** (5%).

### Metrics explained

| Metric | Formula | Meaning |
|---|---|---|
| **Precision@K** | hits@K / K | Of top-K results, what fraction are in the same sub-category as the query? |
| **Recall@K** | hits@K / 249 | Of all 249 same-category products, what fraction appear in top-K? |
| **HitRate@K** | 1 if hits@K ≥ 1 | Did at least one same-category product appear in top-K? |
| **Mean Cosine Similarity** | avg(sim scores of top-K) | How geometrically close are the returned products to the query in embedding space? |

> **Mean Cosine Similarity** is a self-consistency metric, not a quality metric. A high score means the returned products are close to the query vector, but does not guarantee they are useful to the user. Always read it alongside Precision and HitRate.

### Results

**Embedding model:** `text-embedding-3-small` · **Ground truth:** same `sub_category` · 5,000 products evaluated

| K | Precision@K | Recall@K | HitRate@K | Mean Cosine Sim |
|--:|--:|--:|--:|--:|
| 5 | 0.9942 | 0.0200 | 1.0000 | 0.9444 |
| **10** | **0.9956** | **0.0400** | **1.0000** | **0.9103** |
| 20 | 0.9925 | 0.0797 | 1.0000 | 0.8804 |

**Key takeaways:**

- **HitRate@K = 1.0 at all K** — every single product gets at least one same-category result in its top-K. The embedding space is perfectly clean for this dataset.
- **Precision@10 = 0.9956** — nearly all (99.6%) of the top-10 returned products belong to the same sub-category as the query. This is **20× better than random** (5% baseline).
- **Recall stays low** — this is expected and not a flaw. With 249 relevant items per category, a top-10 list can only recover 10/249 ≈ 4% by definition. Recall would only matter if the system were trying to surface all relevant items (e.g. search), not select the best few (recommendation).
- **Mean Cosine Sim ≈ 0.91 at K=10** — the top-10 results sit very close to the query in embedding space, confirming a well-separated latent space.

### How to run (CB)

```bash
cd /path/to/ai-service

# First run: calls OpenAI API, caches embeddings
.venv/bin/python -m app.evaluation.cb_evaluator

# Subsequent runs: loads from cache, no API calls
.venv/bin/python -m app.evaluation.cb_evaluator

# Force re-embed (e.g. after changing the dataset)
.venv/bin/python -m app.evaluation.cb_evaluator --no-cache

# Custom K values
.venv/bin/python -m app.evaluation.cb_evaluator --k-values 5,10,20,50

# Interactive notebook
.venv/bin/jupyter notebook cb_evaluation.ipynb
```

**Source:** `app/evaluation/cb_evaluator.py`

---

## Comparison & Conclusions

| | CF (MovieLens 100k) | CB (ecommerce_products_5k) |
|---|---|---|
| **What is being evaluated** | Recommendation quality for a user based on interaction history | Retrieval quality for a product query based on semantic similarity |
| **Precision@10** | 0.2491 | 0.9956 |
| **HitRate@10** | 0.8134 | 1.0000 |
| **Random baseline P@10** | ~0.024 (2.4%) | ~0.050 (5%) |
| **Lift over random** | ~10× | ~20× |

**Why CF scores lower than CB — and why that is normal:**

The CF and CB systems solve fundamentally different problems and cannot be directly compared.

- **CB** answers: *"Which products look similar to this one?"* — a taxonomy-based retrieval task. The embedding model (`text-embedding-3-small`) was trained on billions of web documents and understands product categories extremely well. Near-perfect Precision@10 is expected.

- **CF** answers: *"Which movies would this user like next, based only on past behavior?"* — a true preference prediction task with no semantic features. The system must infer user taste purely from co-occurrence patterns in interaction data. Achieving Precision@10 = 0.25 (10× random) with a simple cosine similarity approach and no neural embeddings is a solid result.

**These two signals are complementary:**

```
User visits a product page
        │
        ├─► CF: "Users who liked this also liked..." (social proof)
        └─► CB: "Similar products in the same category..." (content discovery)
```

Both can and should be combined in a hybrid system for best results.

---

## Limitations

### CF evaluation limitations

| Limitation | Impact |
|---|---|
| Random 80/20 split (not time-ordered) | May leak future information into training; real-world degradation would be higher |
| Only users with ≥1 relevant test item are evaluated | Cold-start users (few interactions) are excluded |
| No incremental update testing | Batch CF goes stale between runs; not measured here |
| MovieLens is movie data, not e-commerce | Score semantics differ (star rating vs. implicit action score) |

### CB evaluation limitations

| Limitation | Impact |
|---|---|
| Synthetic dataset | Descriptions within a category follow similar templates, which may artificially inflate Precision scores |
| `sub_category = relevant` is a proxy | In production, user preference matters more than taxonomy; a user may dislike most Smartphones |
| Evaluation does not test user-vector quality | The production CB system builds a user vector from interactions before querying Qdrant; that path is not evaluated here |
| No diversity or novelty metrics | A system can score perfect Precision@10 while always returning the same popular items |

---

## File Reference

| File | Purpose |
|---|---|
| `app/evaluation/cf_evaluator.py` | CF evaluation script (CLI + importable functions) |
| `app/evaluation/cb_evaluator.py` | CB evaluation script (CLI + importable functions) |
| `cf_evaluation.ipynb` | Interactive CF notebook with charts and sensitivity analysis |
| `cb_evaluation.ipynb` | Interactive CB notebook with charts and per-category breakdown |
| `ml-100k/` | MovieLens 100k dataset (ratings + 5-fold splits) |
| `data-contentbased/ecommerce_products_5k.csv` | E-commerce product dataset |
| `data-contentbased/structure.md` | Column reference for the product dataset |
| `.cache/cb_eval/embeddings.npy` | Cached OpenAI embeddings (generated on first run) |

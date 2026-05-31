# Product Service — API Reference for External AI Chatbot

**Base URL (Docker Compose / local):** `http://localhost:81`  
**Service prefix:** `/api/product`  
**Full base:** `http://localhost:81/api/product`

> All requests and responses use JSON unless noted.  
> Prices are in **integer (VND or smallest currency unit)** — no decimals.  
> All IDs are **UUID strings**.

---

## Authentication

All endpoints listed here are **public** — no `Authorization` header or login is required.  
Some endpoints optionally accept `X-User-Id` (string, user UUID) to enable personalization (view tracking, search history). Pass it if you have it; omit it if you don't.

---

## Data Models

### Product

```json
{
  "id": "string (UUID)",
  "name": "string",
  "description": "string",
  "images": [
    { "id": "string", "url": "string", "order": 0 }
  ],
  "status": "available | out_of_stock | draft",
  "seller_id": "string (UUID)",
  "rating": 4.5,
  "rate_count": 120,
  "sold_count": 300,
  "is_active": true,
  "is_disabled": false,
  "disable_reason": "string (omitted if not disabled)",
  "price": {
    "min": 50000,
    "max": 150000
  },
  "stock": 50,
  "option_groups": [
    { "key": "color", "values": ["red", "blue"] },
    { "key": "size", "values": ["S", "M", "L"] }
  ],
  "variants": [
    {
      "id": "string (UUID)",
      "sku": "string",
      "options": { "color": "red", "size": "M" },
      "price": 80000,
      "stock": 10,
      "image": "string (URL)",
      "sold_count": 50
    }
  ],
  "category_ids": ["string", "string"],
  "seller_category_ids": ["string"],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Category

```json
{
  "id": "string (UUID)",
  "name": "string",
  "image": "string (URL)",
  "product_count": 42
}
```

### Rating

```json
{
  "id": "string (UUID)",
  "product_id": "string (UUID)",
  "variant_id": "string (UUID)",
  "user": {
    "id": "string",
    "name": "string",
    "email": "string",
    "image": "string (URL)",
    "phone": "string"
  },
  "star": 5,
  "content": "string (optional)",
  "images": [
    { "id": "string", "url": "string" }
  ],
  "rating_response": [
    { "id": "string", "content": "string" }
  ],
  "created_at": "2024-01-01T00:00:00Z",
  "updated_at": "2024-01-01T00:00:00Z"
}
```

### Pagination (used in paginated responses)

```json
{
  "current_page": 1,
  "total_pages": 10,
  "total_items": 95,
  "items_per_page": 10
}
```

---

## Endpoints

---

### 1. Get Product by ID

Fetch a single product with full details including variants, images, and options.

```
GET /api/product/public/:id
```

**Path parameters:**

| Name | Type   | Description   |
|------|--------|---------------|
| `id` | string | Product UUID  |

**Optional headers:**

| Header      | Description                                 |
|-------------|---------------------------------------------|
| `X-User-Id` | User UUID. If provided, records a view event |

**Response `200`:**

```json
{
  "id": "...",
  "name": "Áo thun basic",
  "description": "...",
  "price": { "min": 80000, "max": 120000 },
  "stock": 45,
  "rating": 4.3,
  "rate_count": 80,
  "sold_count": 200,
  "variants": [ ... ],
  "option_groups": [ ... ],
  "category_ids": [ ... ],
  ...
}
```

Returns the full `Product` object.

---

### 3. Search Products

Full-text and filtered search across the product catalog with pagination.

```
GET /api/product/public/search
```

**Query parameters:**

| Name             | Type    | Default | Description                                              |
|------------------|---------|---------|----------------------------------------------------------|
| `search_query`   | string  | —       | Free-text search term (product name, description)        |
| `page`           | int     | 1       | Page number                                              |
| `limit`          | int     | 10      | Items per page                                           |
| `min_price`      | int     | —       | Minimum price filter (inclusive)                         |
| `max_price`      | int     | —       | Maximum price filter (inclusive)                         |
| `min_rating`     | float   | —       | Minimum average rating (e.g. `4.0`)                      |
| `max_rating`     | float   | —       | Maximum average rating                                   |
| `category_ids`   | string  | —       | Comma-separated category UUIDs e.g. `id1,id2`            |
| `sort_by`        | string  | —       | `rating` \| `sold_count` \| `price`                      |
| `sort_direction` | string  | `asc`   | `asc` \| `desc`                                          |

**Optional headers:**

| Header      | Description                               |
|-------------|-------------------------------------------|
| `X-User-Id` | If provided, saves query to search history |

**Response `200`:**

```json
{
  "products": [ { ...product }, ... ],
  "pagination": {
    "current_page": 1,
    "total_pages": 5,
    "total_items": 48,
    "items_per_page": 10
  }
}
```

**Example — search for red shirts under 200,000, sorted by best rating:**

```
GET /api/product/public/search?search_query=áo%20thun&max_price=200000&sort_by=rating&sort_direction=desc
```

---

### 4. Get Products by Seller

List all products from a specific seller, with pagination and filters.

```
GET /api/product/public/products/seller/:sellerId
```

**Path parameters:**

| Name       | Type   | Description  |
|------------|--------|--------------|
| `sellerId` | string | Seller UUID  |

**Query parameters:**

| Name              | Type   | Default | Description                                    |
|-------------------|--------|---------|------------------------------------------------|
| `page`            | int    | 1       | Page number                                    |
| `limit`           | int    | 10      | Items per page                                 |
| `category`        | string | —       | Filter by platform category ID                 |
| `seller_category` | string | —       | Filter by seller's own category ID             |
| `status`          | string | —       | `available` \| `out_of_stock` \| `draft`       |
| `search`          | string | —       | Free-text filter on product name               |
| `sort_by`         | string | —       | Field to sort by                               |
| `sort_direction`  | string | `asc`   | `asc` \| `desc`                                |

**Response `200`:**

```json
{
  "products": [ { ...product }, ... ],
  "pagination": { ... }
}
```

---

### 5. Get Variants by IDs (Batch)

Fetch variant details for a list of variant IDs. Use this when you know specific variant IDs (e.g. from a cart or order) and need full product context.

```
POST /api/product/public/variants/batch
```

**Request body:**

```json
{
  "variant_ids": ["uuid1", "uuid2", "uuid3"]
}
```

| Field        | Type            | Required | Description                |
|--------------|-----------------|----------|----------------------------|
| `variant_ids`| array of string | Yes      | At least 1 variant UUID    |

**Response `200`:**

```json
{
  "variants": [
    {
      "product_name": "Áo thun basic",
      "category_ids": ["cat-uuid-1"],
      "seller_id": "seller-uuid",
      "seller_category_ids": ["sc-uuid-1"],
      "variant": {
        "id": "variant-uuid",
        "sku": "SKU-001",
        "options": { "color": "red", "size": "M" },
        "price": 80000,
        "stock": 10,
        "image": "https://...",
        "sold_count": 30
      }
    }
  ]
}
```

---

### 6. List All Categories

Fetch all product categories available on the platform.

```
GET /api/product/public/category
```

**Query parameters:**

| Name   | Type   | Description                            |
|--------|--------|----------------------------------------|
| `name` | string | Optional name filter / search term     |

**Response `200`:** Array of `Category` objects.

```json
[
  { "id": "...", "name": "Thời trang", "image": "https://...", "product_count": 120 },
  { "id": "...", "name": "Điện tử",    "image": "https://...", "product_count": 80  }
]
```

---

### 7. Get Category by ID

```
GET /api/product/public/category/:id
```

**Path parameters:**

| Name | Type   | Description    |
|------|--------|----------------|
| `id` | string | Category UUID  |

**Response `200`:** Single `Category` object.

---

### 8. Get Seller Categories

List the custom categories a seller has created for organizing their shop.

```
GET /api/product/public/seller/:seller_id/category
```

**Path parameters:**

| Name        | Type   | Description  |
|-------------|--------|--------------|
| `seller_id` | string | Seller UUID  |

**Response `200`:** Array of seller category objects.

```json
[
  { "id": "...", "seller_id": "...", "name": "Áo", "product_count": 15 },
  { "id": "...", "seller_id": "...", "name": "Quần", "product_count": 10 }
]
```

---

### 9. Get Ratings for a Product

Fetch paginated reviews for a specific product, with optional star and image filters.

```
GET /api/product/public/rating/product/:productId
```

**Path parameters:**

| Name        | Type   | Description   |
|-------------|--------|---------------|
| `productId` | string | Product UUID  |

**Query parameters:**

| Name       | Type   | Default | Description                                         |
|------------|--------|---------|-----------------------------------------------------|
| `page`     | int    | 1       | Page number                                         |
| `limit`    | int    | 10      | Items per page                                      |
| `star`     | int    | —       | Filter by star rating: `1` `2` `3` `4` `5`          |
| `hasImage` | string | —       | `true` or `1` — only return ratings with images     |

**Response `200`:**

```json
{
  "ratings": [
    {
      "id": "...",
      "product_id": "...",
      "variant_id": "...",
      "user": { "id": "...", "name": "Nguyễn A", "image": "..." },
      "star": 5,
      "content": "Sản phẩm rất tốt!",
      "images": [ { "id": "...", "url": "https://..." } ],
      "rating_response": [ { "id": "...", "content": "Cảm ơn bạn!" } ],
      "created_at": "2024-01-15T10:00:00Z",
      "updated_at": "2024-01-15T10:00:00Z"
    }
  ],
  "total": 80,
  "page": 1,
  "limit": 10
}
```

---

### 10. Get Vouchers by Seller

List all active vouchers for a seller's shop. Useful for showing discounts a customer can apply.

```
GET /api/product/public/vouchers/seller/:sellerId
```

**Path parameters:**

| Name       | Type   | Description  |
|------------|--------|--------------|
| `sellerId` | string | Seller UUID  |

**Response `200`:** Array of voucher objects.

```json
[
  {
    "id": "...",
    "code": "SALE10",
    "name": "Giảm 10%",
    "description": "...",
    "discount_type": "PERCENTAGE",
    "discount_value": 10,
    "max_discount_value": 50000,
    "min_order_value": 100000,
    "apply_scope": "ALL",
    "total_quantity": 100,
    "usage_limit_per_user": 1,
    "start_time": "2024-01-01T00:00:00Z",
    "end_time": "2024-12-31T23:59:59Z",
    "status": "ACTIVE"
  }
]
```

**Voucher field reference:**

| Field                   | Description                                          |
|-------------------------|------------------------------------------------------|
| `discount_type`         | `FIXED` (flat amount off) \| `PERCENTAGE` (% off)    |
| `discount_value`        | Amount or percentage to deduct                       |
| `max_discount_value`    | Cap on discount for `PERCENTAGE` type                |
| `min_order_value`       | Minimum cart total to apply                          |
| `apply_scope`           | `ALL` (entire shop) \| `CATEGORY` (specific categories) |
| `apply_seller_category_ids` | Category IDs if scope is `CATEGORY`             |

---

## Error Response Format

All errors follow this shape:

```json
{
  "code": "ERROR_CODE_STRING",
  "message": "Human-readable description"
}
```

Common HTTP status codes:

| Status | Meaning                                |
|--------|----------------------------------------|
| 400    | Bad request / validation failed        |
| 404    | Resource not found                     |
| 500    | Internal server error                  |

---

## Quick Reference Table

| Use case                              | Method | Path                                           |
|---------------------------------------|--------|------------------------------------------------|
| Get one product (full detail)         | GET    | `/api/product/public/:id`                      |
| List all products                     | GET    | `/api/product/public`                          |
| Search / filter products              | GET    | `/api/product/public/search`                   |
| Products by seller                    | GET    | `/api/product/public/products/seller/:sellerId`|
| Variant details (batch by IDs)        | POST   | `/api/product/public/variants/batch`           |
| All categories                        | GET    | `/api/product/public/category`                 |
| Category by ID                        | GET    | `/api/product/public/category/:id`             |
| Seller's custom categories            | GET    | `/api/product/public/seller/:seller_id/category`|
| Reviews for a product                 | GET    | `/api/product/public/rating/product/:productId`|
| Vouchers for a seller                 | GET    | `/api/product/public/vouchers/seller/:sellerId`|
| Ask a question (built-in chatbot)     | POST   | `/api/product/public/chatbot/ask`              |

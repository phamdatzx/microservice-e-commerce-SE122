# Order Service — API Reference for External AI Chatbot

**Base URL (Docker Compose / local):** `http://localhost:81`  
**Service prefix:** `/api/order`  
**Full base:** `http://localhost:81/api/order`

> All requests and responses use JSON.  
> Prices / totals are numbers (float for totals, int for item prices — VND or smallest currency unit).  
> All IDs are **UUID strings**.

---

## Authentication

These endpoints rely on **trusted gateway headers** — no `Authorization` token is needed.  
The API gateway (Traefik) injects them after verifying the user's session.

| Header        | Type   | Description                                              |
|---------------|--------|----------------------------------------------------------|
| `X-User-Id`   | string | UUID of the authenticated user. Required on most routes. |
| `X-User-Role` | string | Role string: `customer` \| `seller` \| `admin`. Required on guarded routes. |

> For a chatbot acting on behalf of a customer, pass the customer's UUID as `X-User-Id` and `customer` as `X-User-Role`.

---

## Endpoints Useful for an AI Chatbot

| # | Use Case | Method | Path |
|---|----------|--------|------|
| 1 | List a user's orders (track orders, check status) | GET | `/api/order` |
| 2 | View cart contents | GET | `/api/order/cart` |
| 3 | Count items in cart | GET | `/api/order/cart/count` |
| 4 | Check if a user has purchased a product | POST | `/api/order/verify-purchase` |

---

## Data Models

### Order

```json
{
  "id": "string (UUID)",
  "status": "pending | confirmed | shipping | delivered | cancelled | returned",
  "payment_method": "COD | STRIPE",
  "payment_status": "PENDING | PAID | FAILED",
  "total": 250000.0,
  "phone": "0901234567",
  "delivery_code": "GHN_TRACKING_CODE",
  "item_count": 2,
  "is_rated": false,
  "is_reported": false,
  "user": {
    "id": "string",
    "username": "string",
    "name": "string",
    "email": "string"
  },
  "seller": {
    "id": "string",
    "username": "string",
    "name": "string",
    "email": "string"
  },
  "items": [
    {
      "product_id": "string",
      "variant_id": "string",
      "product_name": "Áo thun basic",
      "variant_name": "Size: M, Color: Red",
      "sku": "SKU-001",
      "price": 80000,
      "image": "https://...",
      "quantity": 2
    }
  ],
  "shipping_address": {
    "full_name": "Nguyễn Văn A",
    "phone": "0901234567",
    "address_line": "123 Đường ABC",
    "ward": "Phường 1",
    "district": "Quận 1",
    "province": "TP. Hồ Chí Minh",
    "country": "Vietnam",
    "latitude": 10.762622,
    "longitude": 106.660172
  },
  "voucher": {
    "code": "SALE10",
    "discount_type": "PERCENTAGE",
    "discount_value": 10
  },
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

**Order status values:**

| Status      | Meaning                             |
|-------------|-------------------------------------|
| `pending`   | Just placed, awaiting confirmation  |
| `confirmed` | Seller confirmed                    |
| `shipping`  | In transit                          |
| `delivered` | Delivered to customer               |
| `cancelled` | Cancelled                           |
| `returned`  | Returned                            |

### Cart Item

```json
{
  "id": "string (cart item UUID)",
  "user_id": "string",
  "seller": {
    "id": "string",
    "name": "string",
    "username": "string",
    "image": "https://...",
    "phone": "string"
  },
  "product": {
    "id": "string",
    "name": "Áo thun basic",
    "seller_id": "string",
    "seller_category_ids": ["string"]
  },
  "variant_id": "string",
  "quantity": 2,
  "variant": {
    "id": "string",
    "sku": "SKU-001",
    "options": { "color": "red", "size": "M" },
    "price": 80000,
    "stock": 10,
    "image": "https://..."
  },
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

---

## Endpoint Details

---

### 1. List User's Orders

Fetch all orders for the authenticated user, with optional status filter and pagination. Use this to help a user track their orders, look up a specific order, or summarize order history.

```
GET /api/order
```

**Required headers:**

| Header      | Value          |
|-------------|----------------|
| `X-User-Id` | User's UUID    |

**Query parameters:**

| Name         | Type   | Default | Description                                           |
|--------------|--------|---------|-------------------------------------------------------|
| `status`     | string | —       | Filter by order status (see status values above)      |
| `page`       | int    | 1       | Page number                                           |
| `limit`      | int    | 10      | Items per page (max 100)                              |
| `sort_by`    | string | —       | `total` \| `created_at`                               |
| `sort_order` | string | `desc`  | `asc` \| `desc`                                       |

**Response `200`:**

```json
{
  "orders": [ { ...order }, { ...order } ],
  "total_count": 25,
  "page": 1,
  "limit": 10,
  "total_pages": 3
}
```

**Example use cases:**

- *"Show me my recent orders"* → `GET /api/order?limit=5`
- *"Do I have any pending orders?"* → `GET /api/order?status=pending`
- *"Has my order been shipped?"* → `GET /api/order` then check `status` and `delivery_code` fields
- *"What did I buy last month?"* → `GET /api/order?sort_by=created_at&sort_order=desc`

---

### 2. View Cart Contents

Fetch the full cart for the authenticated customer, including product names, variant options, prices, and seller info.

```
GET /api/order/cart
```

**Required headers:**

| Header        | Value          |
|---------------|----------------|
| `X-User-Id`   | User's UUID    |
| `X-User-Role` | `customer`     |

**No query parameters.**

**Response `200`:**

```json
{
  "cart_items": [
    {
      "id": "cart-item-uuid",
      "user_id": "user-uuid",
      "seller": { "id": "...", "name": "Shop ABC", "username": "...", "image": "...", "phone": "..." },
      "product": { "id": "...", "name": "Áo thun basic", "seller_id": "...", "seller_category_ids": [] },
      "variant_id": "variant-uuid",
      "quantity": 2,
      "variant": {
        "id": "variant-uuid",
        "sku": "SKU-001",
        "options": { "color": "red", "size": "M" },
        "price": 80000,
        "stock": 10,
        "image": "https://..."
      },
      "created_at": "...",
      "updated_at": "..."
    }
  ]
}
```

**Example use cases:**

- *"What's in my cart?"* → show cart items with product names, variants, quantities, prices
- *"How much will my cart cost?"* → sum `variant.price × quantity` across all items
- *"Is [product] already in my cart?"* → scan `product.name` or `product.id`

---

### 3. Count Cart Items

Get a simple integer count of items currently in the cart.

```
GET /api/order/cart/count
```

**Required headers:**

| Header        | Value       |
|---------------|-------------|
| `X-User-Id`   | User's UUID |
| `X-User-Role` | `customer`  |

**Response `200`:**

```json
{
  "count": 3
}
```

**Example use cases:**

- *"How many items are in my cart?"* → read `count`
- *"Is my cart empty?"* → check if `count` is 0

---

### 4. Verify Purchase

Check whether a specific user has successfully purchased a specific product variant. Useful before allowing actions that require a verified purchase (writing a review, filing a report, etc.).

```
POST /api/order/verify-purchase
```

**No authentication headers required** (public endpoint).

**Request body:**

```json
{
  "user_id": "user-uuid",
  "product_id": "product-uuid",
  "variant_id": "variant-uuid"
}
```

| Field        | Type   | Required | Description          |
|--------------|--------|----------|----------------------|
| `user_id`    | string | Yes      | User UUID            |
| `product_id` | string | Yes      | Product UUID         |
| `variant_id` | string | Yes      | Specific variant UUID|

**Response `200`:**

```json
{
  "has_purchased": true
}
```

**Example use cases:**

- *"Can I write a review for this product?"* → verify purchase first, then allow rating if `has_purchased` is `true`
- *"Has this user bought this item?"* → simple eligibility check

---

## Error Response Format

```json
{
  "code": "ERROR_CODE_STRING",
  "message": "Human-readable description"
}
```

| Status | Meaning                        |
|--------|--------------------------------|
| 400    | Bad request / missing fields   |
| 401    | Missing authentication header  |
| 403    | Wrong role for this endpoint   |
| 404    | Order or resource not found    |
| 500    | Internal server error          |

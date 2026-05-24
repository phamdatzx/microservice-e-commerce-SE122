# AI Chatbot — Response Format Specification

**Service:** `ai-service`  
**Endpoint:** `POST /api/ai/chat`  
**Last updated:** 2026-05-24

---

## Request

```
POST /api/ai/chat
Content-Type: application/json
```

### Headers (injected by Traefik gateway after auth)

| Header        | Type   | Required | Description                          |
|---------------|--------|----------|--------------------------------------|
| `X-User-Id`   | string | No       | Authenticated user UUID              |
| `X-User-Role` | string | No       | `customer` \| `seller` \| `admin`    |
| `X-Username`  | string | No       | Display username                     |

### Body

```json
{
  "message": "Tôi muốn mua máy giặt",
  "chat_history": [],
  "context": {
    "current_product_id": "uuid-or-null",
    "page": "product_detail"
  }
}
```

| Field          | Type           | Required | Description                                                |
|----------------|----------------|----------|------------------------------------------------------------|
| `message`      | string         | Yes      | The user's text message                                    |
| `chat_history` | array \| null  | No       | Previous conversation turns for multi-turn context         |
| `context`      | object \| null | No       | Extra client-side context (current page, product ID, etc.) |

---

## Response

Every response follows a **fixed JSON shape**. All entity arrays are always present (defaulting to `[]`), so the frontend never needs to check for `null`.

```json
{
  "response": {
    "message": "Dưới đây là một số mẫu máy giặt...",
    "products": [
      { "id": "76d4015a-...", "name": "Máy giặt Panasonic cửa trên 10 kg" }
    ],
    "orders": [
      { "id": "order-uuid", "status": "shipping" }
    ],
    "vouchers": [
      { "code": "SALE10", "name": "Giảm 10%", "seller_id": "seller-uuid" }
    ],
    "categories": [
      { "id": "cat-uuid", "name": "Điện tử" }
    ],
    "sellers": [
      { "id": "seller-uuid", "name": "Shop ABC" }
    ]
  }
}
```

### Field Reference

| Field        | Type     | Always present | Description                                           |
|--------------|----------|----------------|-------------------------------------------------------|
| `message`    | string   | ✅             | Natural-language answer in Vietnamese (may contain markdown). |
| `products`   | array    | ✅ (may be `[]`) | Products referenced in the answer.                   |
| `orders`     | array    | ✅ (may be `[]`) | Orders referenced in the answer.                     |
| `vouchers`   | array    | ✅ (may be `[]`) | Vouchers referenced in the answer.                   |
| `categories` | array    | ✅ (may be `[]`) | Categories referenced in the answer.                 |
| `sellers`    | array    | ✅ (may be `[]`) | Sellers referenced in the answer.                    |

### Entity Schemas

**Product**

| Field  | Type   | Description    |
|--------|--------|----------------|
| `id`   | string | Product UUID   |
| `name` | string | Product name   |

**Order**

| Field    | Type   | Description                                                              |
|----------|--------|--------------------------------------------------------------------------|
| `id`     | string | Order UUID                                                               |
| `status` | string | `pending` \| `confirmed` \| `shipping` \| `delivered` \| `cancelled` \| `returned` |

**Voucher**

| Field       | Type   | Description              |
|-------------|--------|--------------------------|
| `code`      | string | Voucher code             |
| `name`      | string | Human-readable label     |
| `seller_id` | string | UUID of the seller/shop  |

**Category**

| Field  | Type   | Description    |
|--------|--------|----------------|
| `id`   | string | Category UUID  |
| `name` | string | Category name  |

**Seller**

| Field  | Type   | Description   |
|--------|--------|---------------|
| `id`   | string | Seller UUID   |
| `name` | string | Shop name     |

---

## Example Responses

### Product search

```json
{
  "response": {
    "message": "Dưới đây là 2 mẫu máy giặt:\n\n1. **Máy giặt Panasonic** — 3,705,000đ\n2. **Máy giặt LG 9kg** — 6,090,000đ",
    "products": [
      { "id": "76d4015a-0475-4075-af24-60d4d6711626", "name": "Máy giặt Panasonic cửa trên 10 kg NA-F10S10BRV" },
      { "id": "2df944ae-18dd-400f-b0b1-b8f4934ebe67", "name": "Máy giặt LG 9kg lồng ngang cửa trước" }
    ],
    "orders": [],
    "vouchers": [],
    "categories": [],
    "sellers": []
  }
}
```

### Order tracking

```json
{
  "response": {
    "message": "Bạn có 1 đơn hàng đang giao:\n- Đơn #abc123 — mã vận đơn: GHN123",
    "products": [],
    "orders": [
      { "id": "abc123-uuid", "status": "shipping" }
    ],
    "vouchers": [],
    "categories": [],
    "sellers": []
  }
}
```

### General greeting (no entities)

```json
{
  "response": {
    "message": "Xin chào! Tôi có thể giúp gì cho bạn?",
    "products": [],
    "orders": [],
    "vouchers": [],
    "categories": [],
    "sellers": []
  }
}
```

---

## Notes for Frontend Developers

1. **`message` contains markdown** — render with a markdown library (bold, lists, line breaks).
2. **Entity arrays are for linking** — use `products[].id` to create clickable product cards/links. The `message` text already describes them in human-readable form.
3. **All arrays are always present** — no need to check for `null` or `undefined`.
4. **Fallback behavior** — if the AI fails to produce structured output, the API will still return the same shape with all entity arrays empty and the raw text in `message`.

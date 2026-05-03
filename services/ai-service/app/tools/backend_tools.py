"""
Backend tools exposed to the LLM via OpenAI function-calling.

Each tool is decorated with ``@tool`` so LangChain can automatically
generate the JSON schema that GPT-4o uses to decide *when* and *how*
to call each function.

Currently all tools are **mocked** — they print the action details and
return placeholder data.  Replace the bodies with real HTTP / DB calls
once the corresponding microservice endpoints are available.
"""

from __future__ import annotations

import json
import logging
from typing import Optional

from langchain_core.tools import tool

logger = logging.getLogger(__name__)


# ── 1. Get Product By ID ────────────────────────────────────────────────

@tool
def get_product_by_id(product_id: str) -> str:
    """Get all information about a specific product by its ID.

    Use this tool when the user asks for details about a particular product,
    such as its name, price, description, stock, variants, or rating.

    Args:
        product_id: The UUID of the product to look up.
    """
    logger.info("[MOCK] get_product_by_id called with product_id=%s", product_id)
    print(f"[TOOL] get_product_by_id(product_id={product_id!r})")

    mock_result = {
        "product_id": product_id,
        "name": "Mock Product",
        "description": "Đây là sản phẩm mô phỏng.",
        "price": {"min": 100000, "max": 200000},
        "stock": 50,
        "rating": 4.5,
        "rate_count": 120,
        "status": "available",
    }
    return json.dumps(mock_result, ensure_ascii=False)


# ── 2. Get Product Reviews ──────────────────────────────────────────────

@tool
def get_product_reviews(product_id: str) -> str:
    """Get customer reviews of a product for review summary.

    Use this tool when the user asks about reviews, opinions, or feedback
    for a specific product.  Returns a list of recent reviews including
    star ratings and review text.

    Args:
        product_id: The UUID of the product whose reviews to fetch.
    """
    logger.info("[MOCK] get_product_reviews called with product_id=%s", product_id)
    print(f"[TOOL] get_product_reviews(product_id={product_id!r})")

    mock_result = [
        {
            "user_name": "Nguyễn Văn A",
            "star": 5,
            "content": "Sản phẩm rất tốt, đóng gói cẩn thận.",
        },
        {
            "user_name": "Trần Thị B",
            "star": 4,
            "content": "Chất lượng ổn, giao hàng hơi chậm.",
        },
        {
            "user_name": "Lê Văn C",
            "star": 3,
            "content": "Sản phẩm tạm được, không như mô tả lắm.",
        },
    ]
    return json.dumps(mock_result, ensure_ascii=False)


# ── 3. Get My Orders ────────────────────────────────────────────────────

@tool
def get_my_orders(status: Optional[str] = None) -> str:
    """Get the current user's orders, optionally filtered by status.

    Use this tool when the user asks about their orders, shipments,
    or purchase history.

    Args:
        status: Optional order status filter. Accepted values:
                "SHIPPING", "COMPLETE", "PAYMENT_WAITING".
                If not provided, returns orders of all statuses.
    """
    logger.info("[MOCK] get_my_orders called with status=%s", status)
    print(f"[TOOL] get_my_orders(status={status!r})")

    mock_result = [
        {
            "order_id": "order-001",
            "product_name": "Mock Product A",
            "quantity": 2,
            "total_price": 400000,
            "status": status or "SHIPPING",
            "created_at": "2026-04-28T10:00:00Z",
        },
        {
            "order_id": "order-002",
            "product_name": "Mock Product B",
            "quantity": 1,
            "total_price": 150000,
            "status": status or "COMPLETE",
            "created_at": "2026-04-25T15:30:00Z",
        },
    ]
    return json.dumps(mock_result, ensure_ascii=False)


# ── 4. Get Seller Info By ID ────────────────────────────────────────────

@tool
def get_seller_info_by_id(seller_id: str) -> str:
    """Get information about a seller by their ID.

    Use this tool when the user asks about a seller, shop, or store —
    for example their name, rating, or product count.

    Args:
        seller_id: The UUID of the seller to look up.
    """
    logger.info("[MOCK] get_seller_info_by_id called with seller_id=%s", seller_id)
    print(f"[TOOL] get_seller_info_by_id(seller_id={seller_id!r})")

    mock_result = {
        "seller_id": seller_id,
        "name": "Shop Mô Phỏng",
        "rating": 4.8,
        "product_count": 235,
        "follower_count": 1500,
        "joined_at": "2024-01-15T00:00:00Z",
    }
    return json.dumps(mock_result, ensure_ascii=False)


# ── 5. Get My Vouchers ──────────────────────────────────────────────────

@tool
def get_my_vouchers() -> str:
    """Get the current user's available vouchers / discount codes.

    Use this tool when the user asks about their vouchers, coupons,
    discount codes, or promotions they can use.
    """
    logger.info("[MOCK] get_my_vouchers called")
    print("[TOOL] get_my_vouchers()")

    mock_result = [
        {
            "voucher_id": "voucher-001",
            "code": "GIAM20K",
            "discount": 20000,
            "type": "fixed",
            "min_order_value": 100000,
            "expires_at": "2026-06-01T00:00:00Z",
        },
        {
            "voucher_id": "voucher-002",
            "code": "SALE10PCT",
            "discount": 10,
            "type": "percentage",
            "min_order_value": 200000,
            "expires_at": "2026-05-15T00:00:00Z",
        },
    ]
    return json.dumps(mock_result, ensure_ascii=False)


# ── Aggregate list for the agent ─────────────────────────────────────────

ALL_TOOLS = [
    get_product_by_id,
    get_product_reviews,
    get_my_orders,
    get_seller_info_by_id,
    get_my_vouchers,
]

"""
Backend tools exposed to the LLM via OpenAI function-calling.

Each tool is decorated with ``@tool`` so LangChain can automatically
generate the JSON schema that GPT-4o uses to decide *when* and *how*
to call each function.

HTTP client structure
---------------------
- ``_product_get`` / ``_product_post`` — thin helpers that call the
  product-service through the Traefik gateway (``PRODUCT_SERVICE_URL``).
- ``_order_get`` — same pattern for the order-service.
- All helpers raise ``RuntimeError`` on non-2xx responses so the LLM
  receives a clear error string instead of an unhandled exception.

Tool naming convention
----------------------
Functions are snake_case.  The docstring **first line** is what the LLM
sees as the tool description; keep it concise and action-oriented.
"""

from __future__ import annotations

import json
import logging
from typing import Optional

import httpx
from langchain_core.tools import tool

from app.core.config import get_settings

logger = logging.getLogger(__name__)


# ---------------------------------------------------------------------------
# HTTP client helpers
# ---------------------------------------------------------------------------

def _product_get(path: str, params: dict | None = None, user_id: str | None = None) -> dict | list:
    """GET request to the product-service public API."""
    settings = get_settings()
    url = f"{settings.PRODUCT_SERVICE_URL}/api/product/public{path}"
    headers = {"X-User-Id": user_id} if user_id else {}

    logger.debug("product-service GET %s params=%s", url, params)
    resp = httpx.get(url, params=params, headers=headers, timeout=10.0)

    if not resp.is_success:
        raise RuntimeError(
            f"product-service error {resp.status_code} for GET {path}: {resp.text[:300]}"
        )
    return resp.json()


def _product_post(path: str, body: dict) -> dict | list:
    """POST request to the product-service public API."""
    settings = get_settings()
    url = f"{settings.PRODUCT_SERVICE_URL}/api/product/public{path}"

    logger.debug("product-service POST %s body=%s", url, body)
    resp = httpx.post(url, json=body, timeout=10.0)

    if not resp.is_success:
        raise RuntimeError(
            f"product-service error {resp.status_code} for POST {path}: {resp.text[:300]}"
        )
    return resp.json()


def _order_get(path: str, user_id: str, user_role: str = "customer", params: dict | None = None) -> dict | list:
    """GET request to the order-service, injecting the gateway auth headers."""
    settings = get_settings()
    url = f"{settings.ORDER_SERVICE_URL}/api/order{path}"
    headers = {"X-User-Id": user_id, "X-User-Role": user_role}

    logger.debug("order-service GET %s params=%s", url, params)
    resp = httpx.get(url, params=params, headers=headers, timeout=10.0)

    if not resp.is_success:
        raise RuntimeError(
            f"order-service error {resp.status_code} for GET {path}: {resp.text[:300]}"
        )
    return resp.json()


def _order_post(path: str, body: dict, user_id: str | None = None, user_role: str | None = None) -> dict | list:
    """POST request to the order-service."""
    settings = get_settings()
    url = f"{settings.ORDER_SERVICE_URL}/api/order{path}"
    headers: dict = {}
    if user_id:
        headers["X-User-Id"] = user_id
    if user_role:
        headers["X-User-Role"] = user_role

    logger.debug("order-service POST %s body=%s", url, body)
    resp = httpx.post(url, json=body, headers=headers, timeout=10.0)

    if not resp.is_success:
        raise RuntimeError(
            f"order-service error {resp.status_code} for POST {path}: {resp.text[:300]}"
        )
    return resp.json()


# ---------------------------------------------------------------------------
# Tool: 1. Get product by ID
# ---------------------------------------------------------------------------

@tool
def get_product_by_id(product_id: str, user_id: Optional[str] = None) -> str:
    """Fetch full details of a single product by its UUID.

    Use this tool when:
    - The user asks for price, stock, variants, options, images, or rating
      of a specific product they have already identified (by ID or from a
      previous search result).
    - You need authoritative product data before making a recommendation.

    Returns a JSON object with fields:
    - ``name``, ``description``, ``status`` (available / out_of_stock / draft)
    - ``price.min`` / ``price.max`` (integer, VND)
    - ``stock`` (integer)
    - ``rating`` (float 1–5), ``rate_count``, ``sold_count``
    - ``variants`` — list of variants each with ``options`` (color/size/…),
      ``price``, ``stock``
    - ``option_groups`` — available option keys and their values
    - ``category_ids``, ``seller_id``

    Args:
        product_id: UUID of the product to fetch.
        user_id:    Optional UUID of the current user. When provided the
                    product-service records a view event for personalisation.
    """
    try:
        data = _product_get(f"/{product_id}", user_id=user_id)
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching product: {exc}"


# ---------------------------------------------------------------------------
# Tool: 2. Get product reviews
# ---------------------------------------------------------------------------

@tool
def get_product_reviews(
    product_id: str,
    page: int = 1,
    limit: int = 5,
    star: Optional[int] = None,
) -> str:
    """Fetch customer reviews for a product.

    Use this tool when:
    - The user asks what customers think about a product.
    - The user wants a summary of reviews (positive / negative).
    - The user wants to know the average rating or the number of reviews.

    Returns a JSON object with:
    - ``ratings`` — list of reviews, each with ``user.name``, ``star`` (1–5),
      ``content`` (review text), ``images`` (list of image URLs).
    - ``total`` — total number of reviews for this product.
    - ``page`` / ``limit`` — current page info.

    Args:
        product_id: UUID of the product.
        page:       Page number (default 1).
        limit:      Number of reviews per page (default 5; max useful is 10).
        star:       Optional — filter to only reviews with this star count (1–5).
    """
    params: dict = {"page": page, "limit": limit}
    if star is not None:
        params["star"] = star

    try:
        data = _product_get(f"/rating/product/{product_id}", params=params)
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching reviews: {exc}"


# ---------------------------------------------------------------------------
# Tool: 3. Search product catalog
# ---------------------------------------------------------------------------

@tool
def search_product_catalog(
    query: str = "",
    min_price: Optional[int] = None,
    max_price: Optional[int] = None,
    min_rating: Optional[float] = None,
    category_ids: Optional[str] = None,
    sort_by: Optional[str] = None,
    sort_direction: str = "desc",
    page: int = 1,
    limit: int = 10,
) -> str:
    """Search the product catalog with keyword and filters via the product-service.

    Use this tool when:
    - The user is browsing or filtering by price, rating, or category and
      wants a paginated list with scores / sold counts.
    - You need structured metadata (price range, sold_count) rather than
      semantic similarity — use ``search_products`` (RAG) for semantic search
      and this tool for filter-heavy or sort-heavy queries.

    Returns a JSON object with:
    - ``products`` — list of product objects.
    - ``pagination`` — ``current_page``, ``total_pages``, ``total_items``.

    Args:
        query:          Free-text search term (product name / description).
                        May be empty when filtering only.
        min_price:      Minimum price filter in VND (integer, inclusive).
        max_price:      Maximum price filter in VND (integer, inclusive).
        min_rating:     Minimum average rating (float, e.g. 4.0).
        category_ids:   Comma-separated category UUIDs to restrict results,
                        e.g. ``"uuid1,uuid2"``.
        sort_by:        ``"rating"`` | ``"sold_count"`` | ``"price"``.
        sort_direction: ``"asc"`` or ``"desc"`` (default ``"desc"``).
        page:           Page number (default 1).
        limit:          Items per page (default 10).
    """
    params: dict = {"page": page, "limit": limit, "sort_direction": sort_direction}
    if query:
        params["search_query"] = query
    if min_price is not None:
        params["min_price"] = min_price
    if max_price is not None:
        params["max_price"] = max_price
    if min_rating is not None:
        params["min_rating"] = min_rating
    if category_ids:
        params["category_ids"] = category_ids
    if sort_by:
        params["sort_by"] = sort_by

    try:
        data = _product_get("/search", params=params)
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error searching products: {exc}"


# ---------------------------------------------------------------------------
# Tool: 4. Get seller categories
# ---------------------------------------------------------------------------

@tool
def get_seller_categories(seller_id: str) -> str:
    """Get the custom product categories that a seller has created for their shop.

    PREFER this tool over ``get_products_by_seller`` when:
    - The user asks generally about what a seller/shop sells without requesting
      a specific product list (e.g. "cửa hàng này bán gì?", "shop này có những
      loại hàng gì?").
    - You want a quick, high-level overview of a seller's inventory before
      deciding whether to fetch actual products.

    These are seller-defined categories (not platform categories), so they
    directly reflect how the seller organises their own shop.

    Returns a JSON array of seller category objects, each with:
    - ``id``            — seller category UUID (use as ``seller_category`` filter
                         in ``get_products_by_seller`` to drill down).
    - ``name``          — human-readable category name (e.g. "Áo", "Quần").
    - ``product_count`` — number of products the seller has in this category.

    Args:
        seller_id: UUID of the seller whose categories to fetch.
    """
    try:
        data = _product_get(f"/seller/{seller_id}/category")
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching seller categories: {exc}"


# ---------------------------------------------------------------------------
# Tool: 5. Get products by seller
# ---------------------------------------------------------------------------

@tool
def get_products_by_seller(
    seller_id: str,
    page: int = 1,
    limit: int = 10,
    search: Optional[str] = None,
    status: Optional[str] = None,
    sort_by: Optional[str] = None,
    sort_direction: str = "desc",
) -> str:
    """List products from a specific seller's shop with pagination and filters.

    Use this tool when:
    - The user wants to browse or search a seller's actual product list.
    - You already know what category they are interested in (use ``get_seller_categories``
      first for a general overview, then this tool to drill into specific items).

    NOTE: For a general "what does this shop sell?" question, call
    ``get_seller_categories`` first — it is faster and gives a clearer overview.
    Only call this tool when the user wants to see actual products.

    Returns a JSON object with:
    - ``products`` — list of product objects from this seller.
    - ``pagination`` — ``current_page``, ``total_pages``, ``total_items``.

    Args:
        seller_id:      UUID of the seller / shop.
        page:           Page number (default 1).
        limit:          Items per page (default 10).
        search:         Optional keyword to filter within this seller's products.
        status:         ``"available"`` | ``"out_of_stock"`` | ``"draft"``.
        sort_by:        Field to sort by (e.g. ``"price"``, ``"rating"``).
        sort_direction: ``"asc"`` or ``"desc"`` (default ``"desc"``).
    """
    params: dict = {"page": page, "limit": limit, "sort_direction": sort_direction}
    if search:
        params["search"] = search
    if status:
        params["status"] = status
    if sort_by:
        params["sort_by"] = sort_by

    try:
        data = _product_get(f"/products/seller/{seller_id}", params=params)
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching seller products: {exc}"


# ---------------------------------------------------------------------------
# Tool: 5. Get seller vouchers
# ---------------------------------------------------------------------------

@tool
def get_seller_vouchers(seller_id: str) -> str:
    """Get all active discount vouchers / coupon codes for a seller's shop.

    Use this tool when:
    - The user asks about discounts, promotions, or vouchers for a specific
      shop or seller.
    - The user wants to know whether there is a coupon they can apply at
      checkout.

    Returns a JSON array of voucher objects, each with:
    - ``code`` — the voucher code to apply at checkout.
    - ``name`` / ``description`` — human-readable label.
    - ``discount_type`` — ``"FIXED"`` (flat VND deduction) or
      ``"PERCENTAGE"`` (% off).
    - ``discount_value`` — the amount or percentage.
    - ``max_discount_value`` — cap on discount for PERCENTAGE type.
    - ``min_order_value`` — minimum cart total required to use this voucher.
    - ``apply_scope`` — ``"ALL"`` (entire shop) or ``"CATEGORY"`` (specific
      seller categories).
    - ``start_time`` / ``end_time`` — validity period (ISO 8601).
    - ``status`` — ``"ACTIVE"`` means it can be used right now.

    Args:
        seller_id: UUID of the seller whose vouchers to fetch.
    """
    try:
        data = _product_get(f"/vouchers/seller/{seller_id}")
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching vouchers: {exc}"


# ---------------------------------------------------------------------------
# Tool: 6. Get categories
# ---------------------------------------------------------------------------

@tool
def get_categories(name_filter: Optional[str] = None) -> str:
    """List all product categories available on the platform.

    Use this tool when:
    - The user asks which categories or product types are available.
    - You need to resolve a category name to its UUID before calling
      ``search_product_catalog`` with ``category_ids``.

    Returns a JSON array of category objects, each with:
    - ``id`` — UUID to use in ``category_ids`` filter.
    - ``name`` — human-readable category name.
    - ``product_count`` — number of products in this category.

    Args:
        name_filter: Optional keyword to filter categories by name
                     (e.g. ``"thời trang"`` to find fashion categories).
    """
    params: dict = {}
    if name_filter:
        params["name"] = name_filter

    try:
        data = _product_get("/category", params=params)
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching categories: {exc}"


# ---------------------------------------------------------------------------
# Tool: 7. Get variant details (batch)
# ---------------------------------------------------------------------------

@tool
def get_variants_by_ids(variant_ids: list[str]) -> str:
    """Fetch full details for one or more product variants by their UUIDs.

    Use this tool when:
    - You have variant IDs (e.g. from an order) and need the product name,
      options (color, size, …), price, and stock for those variants.
    - The user asks about a specific variant they ordered or have in their
      cart.

    Returns a JSON object with a ``variants`` array, each element containing:
    - ``product_name`` — name of the parent product.
    - ``variant.id``, ``variant.sku``, ``variant.options`` (e.g.
      ``{"color": "red", "size": "M"}``), ``variant.price``,
      ``variant.stock``, ``variant.sold_count``.

    Args:
        variant_ids: List of variant UUIDs (at least 1, max ~50 recommended).
    """
    try:
        data = _product_post("/variants/batch", {"variant_ids": variant_ids})
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching variants: {exc}"


# ---------------------------------------------------------------------------
# Tool: 8. Get my orders
# ---------------------------------------------------------------------------

@tool
def get_my_orders(
    user_id: str,
    status: Optional[str] = None,
    page: int = 1,
    limit: int = 10,
    sort_by: Optional[str] = None,
    sort_order: str = "desc",
) -> str:
    """Fetch the current user's order history from the order-service.

    Use this tool when:
    - The user asks about their orders, shipments, or purchase history.
    - The user wants to check the status of a specific order.
    - The user asks "Has my order been shipped?" or "What did I buy recently?"

    Returns a JSON object with:
    - ``orders`` — list of order objects, each containing:
        - ``id`` — order UUID.
        - ``status`` — one of: ``pending``, ``confirmed``, ``shipping``,
          ``delivered``, ``cancelled``, ``returned``.
        - ``total`` — total amount paid (float, VND).
        - ``item_count`` — number of items in the order.
        - ``items`` — list of items, each with ``product_name``,
          ``variant_name``, ``price``, ``quantity``.
        - ``payment_method`` — ``COD`` or ``STRIPE``.
        - ``payment_status`` — ``PENDING``, ``PAID``, or ``FAILED``.
        - ``delivery_code`` — shipping tracking code (if dispatched).
        - ``voucher`` — applied voucher info (if any).
        - ``created_at`` / ``updated_at`` — ISO 8601 timestamps.
    - ``total_count``, ``page``, ``limit``, ``total_pages`` — pagination info.

    Args:
        user_id:    UUID of the user whose orders to fetch. Required.
        status:     Optional filter — ``"pending"``, ``"confirmed"``,
                    ``"shipping"``, ``"delivered"``, ``"cancelled"``,
                    ``"returned"``.
        page:       Page number (default 1).
        limit:      Items per page (default 10, max 100).
        sort_by:    ``"total"`` | ``"created_at"``.
        sort_order: ``"asc"`` or ``"desc"`` (default ``"desc"`` = newest first).
    """
    params: dict = {"page": page, "limit": limit, "sort_order": sort_order}
    if status:
        params["status"] = status
    if sort_by:
        params["sort_by"] = sort_by

    try:
        data = _order_get("", user_id=user_id, params=params)
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching orders: {exc}"


# ---------------------------------------------------------------------------
# Tool: 9. Get my cart
# ---------------------------------------------------------------------------

@tool
def get_my_cart(user_id: str) -> str:
    """Fetch the current user's shopping cart.

    Use this tool when:
    - The user asks "What's in my cart?"
    - The user wants to know the total cost of their cart.
    - The user asks whether a specific product is already in their cart.

    Returns a JSON object with:
    - ``cart_items`` — list of cart items, each containing:
        - ``product.name`` — product name.
        - ``variant.options`` — selected options (e.g. ``{"color": "red", "size": "M"}``).
        - ``variant.price`` — unit price (integer, VND).
        - ``quantity`` — quantity added.
        - ``seller.name`` — shop name.
    To compute total cost, multiply ``variant.price × quantity`` for each item
    and sum them.

    Args:
        user_id: UUID of the authenticated user. Required.
    """
    try:
        data = _order_get("/cart", user_id=user_id, user_role="customer")
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching cart: {exc}"


# ---------------------------------------------------------------------------
# Tool: 10. Get cart item count
# ---------------------------------------------------------------------------

@tool
def get_cart_count(user_id: str) -> str:
    """Get the number of items currently in the user's shopping cart.

    Use this tool when:
    - The user asks "How many items are in my cart?"
    - The user asks "Is my cart empty?"
    - You only need a count, not the full cart contents (faster than get_my_cart).

    Returns a JSON object ``{"count": N}`` where N is the integer number of
    distinct cart items (not total quantity).

    Args:
        user_id: UUID of the authenticated user. Required.
    """
    try:
        data = _order_get("/cart/count", user_id=user_id, user_role="customer")
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error fetching cart count: {exc}"


# ---------------------------------------------------------------------------
# Tool: 11. Verify purchase
# ---------------------------------------------------------------------------

@tool
def verify_purchase(user_id: str, product_id: str, variant_id: str) -> str:
    """Check whether a user has successfully purchased a specific product variant.

    Use this tool when:
    - The user asks "Can I write a review for this product?" — a purchase
      must be verified before they are allowed to rate it.
    - You need to confirm purchase eligibility before enabling an action.

    Returns a JSON object ``{"has_purchased": true | false}``.
    - ``true`` — the user bought this variant and is eligible to review it.
    - ``false`` — no completed order found for this variant.

    Args:
        user_id:    UUID of the user to check.
        product_id: UUID of the product.
        variant_id: UUID of the specific product variant.
    """
    try:
        data = _order_post(
            "/verify-purchase",
            body={"user_id": user_id, "product_id": product_id, "variant_id": variant_id},
        )
        return json.dumps(data, ensure_ascii=False)
    except RuntimeError as exc:
        return f"Error verifying purchase: {exc}"


# ---------------------------------------------------------------------------
# Aggregate list for the agent
# ---------------------------------------------------------------------------

ALL_TOOLS = [
    get_product_by_id,
    get_product_reviews,
    search_product_catalog,
    get_seller_categories,       # overview of a seller's shop — call before get_products_by_seller
    get_products_by_seller,
    get_seller_vouchers,
    get_categories,
    get_variants_by_ids,
    get_my_orders,
    get_my_cart,
    get_cart_count,
    verify_purchase,
]

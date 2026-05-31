from __future__ import annotations

import json
import logging
import re
from typing import Any, Optional

from fastapi import APIRouter, Header
from pydantic import BaseModel, Field

from app.agent.agent import chat

logger = logging.getLogger(__name__)
router = APIRouter()

# ---------------------------------------------------------------------------
# Default (empty) structured response — used as fallback template
# ---------------------------------------------------------------------------

_EMPTY_ENTITIES: dict[str, list] = {
    "products": [],
    "orders": [],
    "vouchers": [],
    "categories": [],
    "sellers": [],
}


# ---------------------------------------------------------------------------
# Request schema
# ---------------------------------------------------------------------------

class ChatMessage(BaseModel):
    """A single turn in the conversation history."""
    role: str = Field(..., description="'user' or 'assistant'")
    content: str = Field(..., description="Message text")


class ChatContext(BaseModel):
    """Extra client-side context forwarded by the frontend."""
    current_product_id: Optional[str] = Field(
        default=None,
        description="Product ID the user is currently viewing on the product detail page.",
    )
    compare_product_ids: Optional[list[str]] = Field(
        default=None,
        description="List of product IDs the user is comparing.",
    )
    current_seller_id: Optional[str] = Field(
        default=None,
        description="Seller ID of the shop page the user is currently browsing.",
    )
    current_order_id: Optional[str] = Field(
        default=None,
        description="Order ID the user is currently viewing.",
    )
    page: Optional[str] = Field(
        default=None,
        description="Current page name, e.g. 'product_detail', 'cart', 'order_list', 'home'.",
    )

    class Config:
        extra = "allow"  # allow additional unknown fields from the frontend


class ChatRequest(BaseModel):
    """Request body for the chat endpoint."""
    message: str
    chat_history: Optional[list[ChatMessage]] = None
    context: Optional[ChatContext] = None


# ---------------------------------------------------------------------------
# Helpers
# ---------------------------------------------------------------------------

def _convert_chat_history(raw: list[ChatMessage] | None) -> list[tuple[str, str]]:
    """
    Convert the frontend's ``[{role, content}, ...]`` format into the
    ``[(role, content), ...]`` tuples that LangChain's
    ``MessagesPlaceholder`` expects.

    LangChain recognises these role strings:
    - ``"human"`` / ``"user"``  → HumanMessage
    - ``"ai"`` / ``"assistant"`` → AIMessage
    """
    if not raw:
        return []

    _ROLE_MAP = {"user": "human", "assistant": "ai"}

    return [
        (_ROLE_MAP.get(msg.role, msg.role), msg.content)
        for msg in raw
    ]


def _build_user_context_block(
    user_id: str | None,
    user_role: str | None,
    username: str | None,
    ctx: ChatContext | None,
) -> str:
    """
    Build a structured context prefix that is injected above the user's
    message so the agent always knows who it is talking to and understands
    the user's current page context.
    """
    lines: list[str] = []

    # ── Auth identity ────────────────────────────────────────────────
    if any([user_id, user_role, username]):
        lines.append("=== THÔNG TIN NGƯỜI DÙNG (do hệ thống cung cấp, đã xác thực) ===")
        if user_id:
            lines.append(f"• user_id   : {user_id}")
        if username:
            lines.append(f"• username  : {username}")
        if user_role:
            lines.append(f"• role      : {user_role}")
        lines.append(
            "→ Khi gọi các tool liên quan đến đơn hàng, giỏ hàng hoặc "
            "bất kỳ tool nào yêu cầu user_id, hãy dùng giá trị trên. "
            "KHÔNG hỏi lại người dùng về user_id."
        )

    # ── Page / product context ───────────────────────────────────────
    if ctx:
        context_lines: list[str] = []

        if ctx.page:
            context_lines.append(f"• Trang hiện tại: {ctx.page}")

        if ctx.current_product_id:
            context_lines.append(
                f"• Người dùng đang xem sản phẩm: {ctx.current_product_id} "
                f"→ nếu họ hỏi về 'sản phẩm này' thì chính là product_id trên."
            )

        if ctx.compare_product_ids:
            ids = ", ".join(ctx.compare_product_ids)
            context_lines.append(
                f"• Người dùng đang so sánh các sản phẩm: [{ids}] "
                f"→ nếu họ yêu cầu so sánh, hãy dùng get_product_by_id cho từng sản phẩm."
            )

        if ctx.current_seller_id:
            context_lines.append(
                f"• Người dùng đang xem shop: {ctx.current_seller_id}"
            )

        if ctx.current_order_id:
            context_lines.append(
                f"• Người dùng đang xem đơn hàng: {ctx.current_order_id}"
            )

        # Include any extra/unknown fields the frontend passes
        known_fields = {
            "current_product_id", "compare_product_ids", "current_seller_id",
            "current_order_id", "page",
        }
        if ctx.__pydantic_extra__:
            for key, value in ctx.__pydantic_extra__.items():
                if key not in known_fields and value is not None:
                    context_lines.append(f"• {key}: {value}")

        if context_lines:
            lines.append("=== NGỮ CẢNH TRANG HIỆN TẠI ===")
            lines.extend(context_lines)

    return "\n".join(lines)


def _parse_agent_response(raw: str) -> dict[str, Any]:
    """
    Parse the agent's raw output into the structured response format.

    The LLM is instructed to return a JSON object with ``message`` +
    entity arrays.  If it fails (returns plain text, wraps in a code
    block, etc.), we gracefully fall back to putting the raw text in
    ``message`` with all entity arrays empty.
    """
    text = raw.strip()

    # Strip markdown code fences (```json ... ``` or ``` ... ```)
    if text.startswith("```"):
        text = re.sub(r"^```(?:json)?\s*", "", text)
        text = re.sub(r"\s*```$", "", text)
        text = text.strip()

    try:
        parsed = json.loads(text)
    except json.JSONDecodeError:
        logger.warning("Agent output is not valid JSON — wrapping as plain message")
        return {"message": raw, **_EMPTY_ENTITIES}

    if not isinstance(parsed, dict) or "message" not in parsed:
        logger.warning("Agent JSON missing 'message' key — wrapping as plain message")
        return {"message": raw, **_EMPTY_ENTITIES}

    # Ensure all entity arrays exist and are lists
    result: dict[str, Any] = {"message": parsed["message"]}
    for key in ("products", "orders", "vouchers", "categories", "sellers"):
        val = parsed.get(key)
        result[key] = val if isinstance(val, list) else []

    return result


# ---------------------------------------------------------------------------
# Endpoint
# ---------------------------------------------------------------------------

@router.post("/chat")
def chat_endpoint(
    body: ChatRequest,
    x_user_id: Optional[str] = Header(default=None, alias="X-User-Id"),
    x_user_role: Optional[str] = Header(default=None, alias="X-User-Role"),
    x_username: Optional[str] = Header(default=None, alias="X-Username"),
) -> dict:
    """
    Main chatbot entry point.

    Headers (injected by the Traefik gateway after auth):
    - ``X-User-Id``   — authenticated user UUID
    - ``X-User-Role`` — user role (customer / seller / admin)
    - ``X-Username``  — username string

    Body fields:
    - ``message``      — the user's text message (required)
    - ``chat_history`` — optional list of ``{role, content}`` objects
                         representing previous turns
    - ``context``      — optional object with client-side context:
        - ``current_product_id`` — product the user is viewing
        - ``compare_product_ids`` — products being compared
        - ``current_seller_id`` — shop page the user is on
        - ``current_order_id`` — order the user is viewing
        - ``page`` — current page name
        - Any additional fields are also forwarded to the agent.
    """
    # Convert frontend chat history to LangChain tuple pairs
    lc_history = _convert_chat_history(body.chat_history)

    # Build the user/session context block that is prepended to the message
    context_block = _build_user_context_block(
        user_id=x_user_id,
        user_role=x_user_role,
        username=x_username,
        ctx=body.context,
    )

    # Compose the full input: context (if any) + separator + user message
    if context_block:
        full_input = f"{context_block}\n\n=== TIN NHẮN CỦA NGƯỜI DÙNG ===\n{body.message}"
    else:
        full_input = body.message

    raw_output = chat(full_input, chat_history=lc_history)

    # Parse the structured JSON from the agent, with graceful fallback
    structured = _parse_agent_response(raw_output)

    return {"response": structured}
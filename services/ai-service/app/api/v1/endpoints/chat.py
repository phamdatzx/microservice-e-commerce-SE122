from __future__ import annotations

import json
import logging
import re
from typing import Any, Optional

from fastapi import APIRouter, Header
from pydantic import BaseModel

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


class ChatRequest(BaseModel):
    """Request body for the chat endpoint."""
    message: str
    chat_history: Optional[list] = None
    # Optional extra context the client may forward (product page context, etc.)
    context: Optional[dict] = None


def _build_user_context_block(
    user_id: str | None,
    user_role: str | None,
    username: str | None,
    body_context: dict | None,
) -> str:
    """
    Build a structured context prefix that is injected above the user's
    message so the agent always knows who it is talking to and can pass
    the correct ``user_id`` to order/cart tools without asking the user.
    """
    lines: list[str] = []

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

    if body_context:
        lines.append("=== NGỮ CẢNH BỔ SUNG TỪ CLIENT ===")
        for key, value in body_context.items():
            lines.append(f"• {key}: {value}")

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
    - ``chat_history`` — optional list of previous turns for multi-turn conversations
    - ``context``      — optional dict with extra client-side context
                         (e.g. ``{"current_product_id": "uuid", "page": "product_detail"}``)
    """
    # Build the user/session context block that is prepended to the message
    context_block = _build_user_context_block(
        user_id=x_user_id,
        user_role=x_user_role,
        username=x_username,
        body_context=body.context,
    )

    # Compose the full input: context (if any) + separator + user message
    if context_block:
        full_input = f"{context_block}\n\n=== TIN NHẮN CỦA NGƯỜI DÙNG ===\n{body.message}"
    else:
        full_input = body.message

    raw_output = chat(full_input, chat_history=body.chat_history)

    # Parse the structured JSON from the agent, with graceful fallback
    structured = _parse_agent_response(raw_output)

    return {"response": structured}
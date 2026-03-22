from datetime import datetime

from pydantic import BaseModel, Field


class UserInteractionMessage(BaseModel):
    """
    RabbitMQ payload for a single user–product interaction.

    If `score` is omitted, the worker resolves it from `action` using
    `INTERACTION_ACTION_SCORES_JSON` (see Settings).
    """

    user_id: str = Field(..., min_length=1)
    product_id: str = Field(..., min_length=1)
    action: str = Field(..., min_length=1)
    score: float | None = Field(
        default=None,
        description="Optional explicit weight; overrides action-based score when set.",
    )
    timestamp: datetime | None = Field(
        default=None,
        description="Optional event time (ISO-8601). Defaults to server UTC now.",
    )

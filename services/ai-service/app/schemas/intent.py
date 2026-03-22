from pydantic import BaseModel, Field

class IntentRequest(BaseModel):
    text: str = Field(..., description="The user's message/text to predict the intent for", example="Cho mình hỏi giá chiếc iPhone này")

class IntentResponse(BaseModel):
    intent: str = Field(..., description="The predicted intent label", example="ask_price")

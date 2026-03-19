from fastapi import APIRouter, HTTPException
from typing import Any
import logging

from app.schemas.intent import IntentRequest, IntentResponse
from app.services.predict_intent_service import predict_intent

router = APIRouter()
logger = logging.getLogger(__name__)

@router.post("/predict", response_model=IntentResponse)
def predict_user_intent(request: IntentRequest) -> Any:
    """
    Predict the intent of the given user text.
    """
    try:
        intent = predict_intent(request.text)
        return {"intent": intent}
    except Exception as e:
        logger.error(f"Error predicting intent: {e}")
        raise HTTPException(status_code=500, detail="Internal server error during intent prediction")

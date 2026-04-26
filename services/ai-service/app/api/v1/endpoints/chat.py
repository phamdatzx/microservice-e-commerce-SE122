from fastapi import APIRouter
from app.services.intent_router import route_intent

router = APIRouter()

@router.post("/chat")
def chat(query: str):
    result = route_intent(query)
    return {"response": result}
from fastapi import APIRouter
from app.services.router import route_intent

router = APIRouter()

@router.post("/chat")
def chat(query: str):
    result = route_intent(query)
    return {"response": result}
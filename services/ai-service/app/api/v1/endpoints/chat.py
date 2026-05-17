from fastapi import APIRouter
from app.agent.agent import chat

router = APIRouter()


@router.post("/chat")
def chat_endpoint(query: str):
    result = chat(query)
    return {"response": result}
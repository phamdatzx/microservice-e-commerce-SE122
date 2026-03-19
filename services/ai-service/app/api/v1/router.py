from fastapi import APIRouter

from app.api.v1.endpoints.health import router as health_router
from app.api.v1.endpoints.product_vectors import router as product_vectors_router
from app.api.v1.endpoints.user_vectors import router as user_vectors_router
from app.api.v1.endpoints.recommendations import router as recommendations_router
from app.api.v1.endpoints.intent import router as intent_router

api_router = APIRouter()
api_router.include_router(health_router, tags=["health"])
api_router.include_router(product_vectors_router, tags=["products"])
api_router.include_router(user_vectors_router, tags=["users"])
api_router.include_router(recommendations_router, tags=["recommendations"])
api_router.include_router(intent_router, prefix="/intent", tags=["intent"])

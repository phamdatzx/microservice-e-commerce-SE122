from fastapi import APIRouter

from app.schemas.health import HealthResponse, RootResponse
from app.services.health_service import get_health_status, get_root_message

router = APIRouter()


@router.get("/", response_model=RootResponse)
async def root() -> RootResponse:
    return get_root_message()


@router.get("/health", response_model=HealthResponse)
async def health_check() -> HealthResponse:
    return get_health_status()

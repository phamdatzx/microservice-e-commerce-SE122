from fastapi import FastAPI

from app.api.v1.router import api_router
from app.core.config import get_settings
from app.services.embedding_service import get_embedding_model


def create_app() -> FastAPI:
    settings = get_settings()

    application = FastAPI(
        title=settings.PROJECT_NAME,
        version=settings.VERSION,
        debug=settings.DEBUG,
    )
    application.include_router(api_router, prefix=settings.API_V1_PREFIX)

    @application.on_event("startup")
    async def load_embedding_model() -> None:  # noqa: D401
        """Warm up embedding model at startup."""
        get_embedding_model()

    return application

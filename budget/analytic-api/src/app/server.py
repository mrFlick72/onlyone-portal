import logging
import os

from app.container import ApplicationContainer
from app.infrastructure.middleware.security_context_injector_filter import SecurityContextInjectorFilter
from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.analytic.api.end_point import analytic_end_point_router
from app.infrastructure.management.health_end_point import health_end_point_router

from dotenv import load_dotenv


load_dotenv(dotenv_path=os.getenv("ANALYTIC_API_CONFIG_FILE_LOCATION"))

logger = logging.getLogger(__name__)
application_container = ApplicationContainer()
application_container.wire(modules=["app.analytic.api.end_point"])

app = FastAPI()
origins = [o for o in os.getenv("CORS_ALLOWED_ORIGINS", "").split(",") if o]
if not origins:
    raise RuntimeError("CORS_ALLOWED_ORIGINS must be set to an explicit origin list — wildcard is not allowed with credentials")
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.state.container = application_container

if os.getenv("WITH_MIDDLEWARE", "true").lower() == "true":
    app.add_middleware(
        SecurityContextInjectorFilter,
        "security_context",
        application_container.security_context_config_container.security_context_resolver(),
    )

logger.info(
    "Middleware loaded: %s",
    application_container.security_context_config_container.security_context_resolver(),
)
app.include_router(health_end_point_router)
app.include_router(analytic_end_point_router)

import logging
import os

from fastapi import FastAPI
from fastapi.middleware.cors import CORSMiddleware

from app.revenue.api.end_point import revenue_end_point_router
from app.infrastructure.management.health_end_point import health_end_point_router
from app.infrastructure.middleware.user_name_injector_filter import (
    UserNameInjectorFilter,
)
from app.container import ApplicationContainer

from dotenv import load_dotenv


load_dotenv(dotenv_path=os.getenv("BUDGET_API_CONFIG_FILE_LOCATION"))

logger = logging.getLogger(__name__)  # noqa: F821
application_container = ApplicationContainer()  # type: ignore
application_container.wire(modules=["app.revenue.api.end_point"])

app = FastAPI()
origins = os.getenv("CORS_ALLOWED_ORIGINS", "*").split(",")
app.add_middleware(
    CORSMiddleware,
    allow_origins=origins,
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)
app.container = application_container

# Set up application middleware
if os.getenv("WITH_MIDDLEWARE", "true").lower() == "true":
    app.add_middleware(
        UserNameInjectorFilter,
        "user_name",
        application_container.user_config_container.user_name_resolver(),
    )

logger.info(
    "Middleware loaded: %s",
    application_container.user_config_container.user_name_resolver(),
)
app.include_router(health_end_point_router)
app.include_router(revenue_end_point_router)

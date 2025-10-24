import os

from fastapi import FastAPI
from app.revenue.api.end_point import revenue_end_point_router
from app.infrastructure.management.health_end_point import health_end_point_router
from app.infrastructure.middleware.user_name_injector_filter import (
    UserNameInjectorFilter,
)
from app.container import ApplicationContainer

from dotenv import load_dotenv

load_dotenv(dotenv_path=os.getenv("BUDGET_API_CONFIG_FILE_LOCATION"))


application_container = ApplicationContainer()  # type: ignore
application_container.wire(modules=["app.revenue.api.end_point"])

app = FastAPI()
app.container = application_container

# Set up application middleware
if os.getenv("WITH_MIDDLEWARE", "true").lower() == "true":
    app.add_middleware(
        UserNameInjectorFilter,
        "user_name",
        application_container.user_config_container.user_name_resolver(),
    )

print("Middleware loaded:", application_container.user_config_container.user_name_resolver())
app.include_router(health_end_point_router)
app.include_router(revenue_end_point_router)

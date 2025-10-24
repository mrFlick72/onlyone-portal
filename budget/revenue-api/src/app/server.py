import os

from fastapi import FastAPI
from app.revenue.api.revenue_end_point import revenue_end_point_router
from app.infrastructure.management.health_end_point import health_end_point_router
from app.infrastructure.middleware.user_name_injector_filter import (
    UserNameInjectorFilter,
)
from app.user.container import UserConfigContainer

from dotenv import load_dotenv

load_dotenv(dotenv_path=os.getenv("BUDGET_API_CONFIG_FILE_LOCATION"))

user_config_container = UserConfigContainer()


app = FastAPI()

# Set up application middleware
if os.getenv("WITH_MIDDLEWARE", "true").lower() == "true":
    app.add_middleware(
        UserNameInjectorFilter,
        "user_name",
        user_config_container.get_user_name_resolver(),
    )

print("Middleware loaded:", user_config_container.get_user_name_resolver())
print("Middleware loaded:", user_config_container.get_user_name_resolver())
app.include_router(health_end_point_router)
app.include_router(revenue_end_point_router)

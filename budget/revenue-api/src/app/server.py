import os

from flask import Flask
from app.revenue.api.revenue_end_point import revenue_end_point
from app.infrastructure.management.health_end_point import HealthEndPoint
from app.infrastructure.middleware.user_name_injector_filter import (
    UserNameInjectorFilter,
)
from app.user.adapter.thread.local_thread_user_name_resolver import (
    LocalThreadUserNameResolver,
)
from dotenv import load_dotenv

load_dotenv(dotenv_path=os.getenv("BUDGET_API_CONFIG_FILE_LOCATION"))

app = Flask(__name__)

# Set up application middleware
if os.getenv("WITH_MIDDLEWARE", "true").lower() == "true":
    user_name_injector_filter = UserNameInjectorFilter(
        LocalThreadUserNameResolver().get_instance()
    )
    app.before_request(user_name_injector_filter.filter)

# Register endpoint routes
HealthEndPoint(app)
app.register_blueprint(revenue_end_point)

import os

from flask import Flask
from app.revenue.api.revenue_end_point import revenue_end_point
from app.infrastructure.management.health_end_point import HealthEndPoint
from app.infrastructure.middleware.user_name_injector_filter import UserNameInjectorFilter
from dotenv import load_dotenv

load_dotenv(dotenv_path=os.getenv("BUDGET_API_CONFIG_FILE_LOCATION"))


app = Flask(__name__)

# user_name_injector_filter = UserNameInjectorFilter()
# app.before_request(user_name_injector_filter.filter)

HealthEndPoint(app)
app.register_blueprint(revenue_end_point)

from flask import Blueprint, request
from app.revenue.config import RevenueConfigurationProvider
from app.revenue.api.revenue_converter import from_representation_to_domain

revenue_end_point = Blueprint("revenue_end_point", __name__)


_NO_CONTENT = ""


@revenue_end_point.route("/budget/revenue", methods=["GET"])
def get_revenue():
    return {}


@revenue_end_point.route("/budget/revenue", methods=["POST"])
def save_revenue():
    save_revenue_service = RevenueConfigurationProvider.get_save_revenue_service()
    revenue_representation = request.get_json()
    revenue = from_representation_to_domain(revenue_representation)
    save_revenue_service.save_revenue(revenue)
    return _NO_CONTENT, 201


@revenue_end_point.route("/budget/revenue/<id>", methods=["PUT"])
def update_revenue():
    return {}


@revenue_end_point.route("/budget/revenue/<id>", methods=["DELETE"])
def delete_revenue():
    return {}

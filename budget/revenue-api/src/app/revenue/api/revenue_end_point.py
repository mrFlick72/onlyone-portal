from flask import Blueprint, request
from app.revenue.config import RevenueConfigurationProvider
from app.revenue.api.revenue_converter import fromRepresentationToDomain

revenue_end_point = Blueprint("revenue_end_point", __name__)

_save_revenue_service = RevenueConfigurationProvider.get_save_revenue_service()

_NO_CONTENT = ""


@revenue_end_point.route("/budget/revenue", methods=["GET"])
def get_revenue():
    return {}


@revenue_end_point.route("/budget/revenue", methods=["POST"])
def save_revenue():
    revenue_representation = request.get_json()
    print(f"revenue_representation: {revenue_representation}")
    revenue = fromRepresentationToDomain(revenue_representation)
    print(f"revenue: {revenue}")
    _save_revenue_service.save_revenue(revenue)
    return _NO_CONTENT, 201


@revenue_end_point.route("/budget/revenue/<id>", methods=["PUT"])
def update_revenue():
    return {}


@revenue_end_point.route("/budget/revenue/<id>", methods=["DELETE"])
def delete_revenue():
    return {}

from fastapi import APIRouter,Response
from app.revenue.config import RevenueConfigurationProvider
from app.revenue.api.revenue_converter import from_representation_to_domain
from app.revenue.api.representation import RevenueRepresentation

revenue_end_point_router = APIRouter()


@revenue_end_point_router.route("/budget/revenue", methods=["GET"])
async def get_revenue():
    return {}


@revenue_end_point_router.post("/budget/revenue")
async def save_revenue(representation: RevenueRepresentation):
    save_revenue_service = RevenueConfigurationProvider.get_save_revenue_service()
    revenue = from_representation_to_domain(representation)
    save_revenue_service.save_revenue(revenue)
    return Response(status_code=201)


@revenue_end_point_router.route("/budget/revenue/{id}", methods=["PUT"])
async def update_revenue():
    return {}


@revenue_end_point_router.route("/budget/revenue/{id}", methods=["DELETE"])
async def delete_revenue():
    return {}

from fastapi import APIRouter, Response,Depends
from dependency_injector.wiring import Provide
from typing import Annotated

from app.revenue.api.representation import RevenueRepresentation
from app.revenue.api.converter import RevenueConverter
from app.user.container import UserConfigContainer
from app.revenue.container import RevenueConfigContainer
from app.user.domain.user_name_resolver import UserNameResolver
from app.revenue.domain.service import SaveRevenue  

revenue_end_point_router = APIRouter()


@revenue_end_point_router.route("/budget/revenue", methods=["GET"])
async def get_revenue():
    return {}


@revenue_end_point_router.post("/budget/revenue")
async def save_revenue(
    representation: RevenueRepresentation,
    user_name_resolver=Annotated[
        UserNameResolver, Depends[Provide[UserConfigContainer.get_user_name_resolver]]
    ],
    save_revenue_service=Annotated[
        SaveRevenue, Depends[Provide[RevenueConfigContainer.save_revenue_service]]
    ],
):
    converter = RevenueConverter(user_name_resolver())

    revenue = converter.from_representation_to_domain(representation)
    save_revenue_service.save_revenue(revenue)
    return Response(status_code=201)


@revenue_end_point_router.route("/budget/revenue/{id}", methods=["PUT"])
async def update_revenue():
    return {}


@revenue_end_point_router.route("/budget/revenue/{id}", methods=["DELETE"])
async def delete_revenue():
    return {}

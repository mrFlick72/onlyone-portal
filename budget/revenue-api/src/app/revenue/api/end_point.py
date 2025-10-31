import json
from fastapi import APIRouter, Response, Depends
from dependency_injector.wiring import Provide, inject
from typing import Annotated

from app.revenue.api.representation import (
    RevenueRequestRepresentation,
)
from app.revenue.api.converter import QueryParamRepresentationConverter
from app.revenue.api.converter import RevenueConverter
from app.revenue.domain.service import SaveRevenue, DeleteRevenue, FindRevenue
from app.container import ApplicationContainer
from app.revenue.domain.revenue import RevenueId

revenue_end_point_router = APIRouter()


@revenue_end_point_router.get("/budget/revenue")
@inject
async def get_revenue(
    q: str,
    find_revenue_service: Annotated[
        FindRevenue,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.find_revenue_service]
        ),
    ],
    query_param_converter: Annotated[
        QueryParamRepresentationConverter,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.query_param_converter]
        ),
    ],
    converter: Annotated[
        RevenueConverter,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.revenue_converter]
        ),
    ],
):
    query_param_representation = query_param_converter.from_query_param_to_year(q)
    revenues = find_revenue_service.find_by(query_param_representation.year)
    print(f"revenues from the database {revenues}")
    json_response = [
        converter.from_domain_to_representation(revenue).model_dump()
        for revenue in revenues
    ]
    print(json.dumps(json_response))
    

    return Response(
        status_code=200,
        content = json.dumps(json_response),
        media_type="application/json",
    )


@revenue_end_point_router.post("/budget/revenue")
@inject
async def save_revenue(
    representation: RevenueRequestRepresentation,
    save_revenue_service: Annotated[
        SaveRevenue,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.save_revenue_service]
        ),
    ],
    converter: Annotated[
        RevenueConverter,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.revenue_converter]
        ),
    ],
):
    revenue = converter.from_representation_to_domain(representation)
    save_revenue_service.save(revenue)
    return Response(status_code=201)


@revenue_end_point_router.put("/budget/revenue/{id}")
@inject
async def update_revenue(
    id: str,
    representation: RevenueRequestRepresentation,
    save_revenue_service: Annotated[
        SaveRevenue,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.save_revenue_service]
        ),
    ],
    converter: Annotated[
        RevenueConverter,
        Depends(
            Provide[ApplicationContainer.revenue_config_container.revenue_converter]
        ),
    ],
):
    revenue = converter.from_representation_to_domain(representation)
    revenue.id = RevenueId(id)
    save_revenue_service.save(revenue)
    return Response(status_code=204)


@revenue_end_point_router.delete("/budget/revenue/{id}")
@inject
async def delete_revenue(
    id: str,
    delete_revenue_service: Annotated[
        DeleteRevenue,
        Depends(
            Provide[
                ApplicationContainer.revenue_config_container.delete_revenue_service
            ]
        ),
    ],
):
    delete_revenue_service.delete(RevenueId(id))
    return Response(status_code=204)

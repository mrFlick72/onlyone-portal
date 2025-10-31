from app.revenue.domain.revenue import Revenue
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.time.domain.year import Year
from app.revenue.api.representation import RevenueRequestRepresentation,RevenueResponseRepresentation
from app.user.domain.user_name_resolver import UserNameResolver
from app.revenue.api.representation import QueryParamRepresentation


class RevenueConverter:

    def __init__(self, user_name_resolver: UserNameResolver):
        self.user_name_resolver = user_name_resolver

    def from_representation_to_domain(
        self, representation: RevenueRequestRepresentation
    ) -> Revenue:
        return Revenue(
            id=None,
            user_name=self.user_name_resolver.get_user_name(),  # type: UserName
            date=Date.date_for(representation.date),
            amount=Money.money_for(representation.amount),
            note=representation.note,
        )

    # to be tested
    def from_domain_to_representation(self, revenue: Revenue) -> RevenueResponseRepresentation:
        return RevenueResponseRepresentation(
            id=revenue.id.content,
            user_name=revenue.user_name.content,  # type: UserName
            date=revenue.date.formatted_date(),
            amount=revenue.amount.stringify_amount(),
            note=revenue.note,
        )


class QueryParamRepresentationConverter:

    def from_query_param_to_year(self, query_param: str) -> QueryParamRepresentation:
        splitted_params = query_param.split(";")
        year_param = splitted_params[0].split("=")[1]
        return QueryParamRepresentation(year=Year(int(year_param)))

from app.revenue.domain.revenue import Revenue
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.revenue.api.representation import RevenueRepresentation
from app.user.domain.user_name_resolver import UserNameResolver


class RevenueConverter:

    def __init__(self, user_name_resolver: UserNameResolver):
        self.user_name_resolver = user_name_resolver

    def from_representation_to_domain(
        self, representation: RevenueRepresentation
    ) -> Revenue:
        return Revenue(
            id=None,
            user_name=self.user_name_resolver.get_user_name(),  # type: UserName
            date=Date.date_for(representation.date),
            amount=Money.money_for(representation.amount),
            note=representation.note,
        )
        

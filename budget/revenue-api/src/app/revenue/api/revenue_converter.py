from app.revenue.domain.revenue import Revenue
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.user.config import get_user_name_resolver
from app.revenue.api.representation import RevenueRepresentation

def from_representation_to_domain(representation: RevenueRepresentation) -> Revenue:
    user_name_resolver = get_user_name_resolver()

    return Revenue(
        id=representation.id,
        user_name=user_name_resolver.get_user_name(),  # type: UserName
        date=Date.date_for(representation.date),
        amount=Money.money_for(representation.amount),
        note=representation.note,
    )

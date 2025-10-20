from app.revenue.domain.revenue import Revenue
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.user.config import get_user_name_resolver


def from_representation_to_domain(representation: dict) -> Revenue:
    user_name_resolver = get_user_name_resolver()

    return Revenue(
        id=None,
        user_name=user_name_resolver.get_user_name(),  # type: UserName
        date=Date.date_for(representation["date"]),
        amount=Money.money_for(representation["amount"]),
        note=representation["note"],
    )

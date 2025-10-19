from app.revenue.domain.revenue import Revenue
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.user.config import user_name_resolver


_user_name_resolver = user_name_resolver()


def fromRepresentationToDomain(representation: dict) -> Revenue:
    print(f"fromRepresentationToDomain called with: {representation}")
    return Revenue(
        id=None,
        user_name=_user_name_resolver.get_user_name(),  # type: UserName
        date=Date.date_for(representation["date"]),
        amount=Money.money_for(representation["amount"]),
        note=representation["note"],
    )

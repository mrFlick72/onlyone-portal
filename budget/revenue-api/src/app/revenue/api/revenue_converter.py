from app.revenue.domain.revenue import Revenue
from app.money.domain.money import Money
from app.time.domain.date import Date
from app.user.domain.user import UserName


def fromRepresentationToDomain(representation: dict) -> Revenue:
    return Revenue(
        id=None,
        user_name=UserName(representation["user_name"]),
        registration_date=Date.date_for(representation["date"]),
        amount=Money.money_for(representation["amount"]),
        note=representation["note"],
    )

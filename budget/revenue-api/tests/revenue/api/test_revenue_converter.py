import pytest

from app.revenue.api.revenue_converter import fromRepresentationToDomain
from app.money.domain.money import Money
from app.time.domain.date import Date, DateParsingException
from app.user.domain.user import UserName


def test_from_representation_to_domain_happy_path():
    rep = {
        "user_name": "alice",
        "registration_date": "15/08/2025",
        "amount": "123.45",
        "note": "salary",
    }

    revenue = fromRepresentationToDomain(rep)

    # id is set to None by the converter
    assert revenue.id is None
    assert isinstance(revenue.user_name, UserName)
    assert revenue.user_name.content == "alice"
    assert isinstance(revenue.registration_date, Date)
    assert revenue.registration_date.formatted_date() == "15/08/2025"
    assert isinstance(revenue.amount, Money)
    assert revenue.amount.stringify_amount() == str(revenue.amount.amount)
    assert revenue.note == "salary"


def test_from_representation_to_domain_invalid_date_raises():
    rep = {
        "user_name": "bob",
        "registration_date": "2025-08-15",
        "amount": "10.00",
        "note": "bonus",
    }

    with pytest.raises(DateParsingException):
        fromRepresentationToDomain(rep)

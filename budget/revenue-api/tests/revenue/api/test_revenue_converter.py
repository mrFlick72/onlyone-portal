import pytest

from app.revenue.api.revenue_converter import from_representation_to_domain
from app.money.domain.money import Money
from app.time.domain.date import Date, DateParsingException
from app.user.domain.user import UserName
from app.user.adapter.thread.local_thread_user_name_resolver import (
    LocalThreadUserNameResolver,
)

instance = LocalThreadUserNameResolver.get_instance()


@pytest.fixture(autouse=True)
def reset_user_name_resolver():
    instance.set_user_name(UserName("alice"))

    yield

    instance.set_user_name(None)


def test_from_representation_to_domain_happy_path():
    rep = {
        "date": "15/08/2025",
        "amount": "123.45",
        "note": "salary",
    }

    revenue = from_representation_to_domain(rep)

    # id is set to None by the converter
    assert revenue.id is None
    assert isinstance(revenue.user_name, UserName)
    assert revenue.user_name.content == "alice"
    assert isinstance(revenue.date, Date)
    assert revenue.date.formatted_date() == "15/08/2025"
    assert isinstance(revenue.amount, Money)
    assert revenue.amount.stringify_amount() == str(revenue.amount.amount)
    assert revenue.note == "salary"


def test_from_representation_to_domain_invalid_date_raises():
    rep = {
        "date": "2025-08-15",
        "amount": "10.00",
        "note": "bonus",
    }

    with pytest.raises(DateParsingException):
        from_representation_to_domain(rep)

from decimal import Decimal, getcontext
from app.money.domain.money import Money

def test_get_money_from_a_string():
    first_amount = Money.money_for("20.50")
    second_amount = Money.money_for("20.51")
    actual = first_amount.plus(second_amount)
    print(actual)
    print(Decimal("41.01"))
    
    assert actual == Decimal("41.01")

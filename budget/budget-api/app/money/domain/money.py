# static final int SCALE_PRECISION = 2;
# static final int SCALE_CRITERIA = BigDecimal.ROUND_HALF_DOWN;

# public static final Money ZERO = Money.moneyFor("0.00");
# public static final Money ONE = Money.moneyFor("1.00");

# public Money(BigDecimal amount) {
#     this.amount = amount.setScale(SCALE_PRECISION, SCALE_CRITERIA);
# }

# public static Money moneyFor(String amount) {
#     return new Money(new BigDecimal(amount));
# }

# public Money plus(Money money) {
#     return new Money(this.amount.add(money.amount()).setScale(SCALE_PRECISION, SCALE_CRITERIA));
# }

# public String stringifyAmount() {
#     return amount.toString();
# }

# @Override
# public boolean equals(Object o) {
#     if (this == o) return true;
#     if (o == null || getClass() != o.getClass()) return false;
#     Money money = (Money) o;
#     return Objects.equals(amount, money.amount);
# }


import decimal
from typing import Final

SCALE_PRECISION: Final = 2
SCALE_CRITERIA: Final = decimal.ROUND_HALF_DOWN


class Money:

    def __init__(self, amount: decimal):
        self.amount = amount

    @staticmethod
    def money_for(amount: str):
        pass

    def stringify_amount():
        pass

    def plus(self, money: Money):
        pass

    def __eq__(self, other):
        return self.amount == other.amount

    def __hash__(self):
        return hash((self.amount))

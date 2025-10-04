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
from decimal import ROUND_HALF_DOWN, Context, Decimal
from typing import Final

SCALE_PRECISION: Final = 2
SCALE_CRITERIA: Final = ROUND_HALF_DOWN

MONEY_ARITHMETIC_CONTEXT: Final = Context(rounding=SCALE_CRITERIA)


class Money:

    def __init__(self, amount: decimal):
        self.amount = amount

    @staticmethod
    def money_for(amount: str):
        context = MONEY_ARITHMETIC_CONTEXT
        return Money(context.create_decimal(Decimal(amount)))

    def stringify_amount(self):
        return str(self.amount)

    def plus(self, addend):
        context = MONEY_ARITHMETIC_CONTEXT

        first_amount = self.amount
        second_amount = addend.amount
        result = round(first_amount + second_amount, SCALE_PRECISION)
        return context.create_decimal(result)

from dataclasses import dataclass
from decimal import ROUND_HALF_DOWN, Context, Decimal
from typing import Final

SCALE_PRECISION: Final = 2
SCALE_CRITERIA: Final = ROUND_HALF_DOWN

MONEY_ARITHMETIC_CONTEXT: Final = Context(rounding=SCALE_CRITERIA)


@dataclass
class Money:
    amount: Decimal

    @staticmethod
    def money_for(amount: str) -> "Money":
        context = MONEY_ARITHMETIC_CONTEXT
        return Money(context.create_decimal(Decimal(amount)))

    def stringify_amount(self) -> str:
        return str(self.amount)

    def plus(self, addend: "Money") -> Decimal:
        context = MONEY_ARITHMETIC_CONTEXT

        first_amount = self.amount
        second_amount = addend.amount
        result = round(first_amount + second_amount, SCALE_PRECISION)
        return context.create_decimal(result)

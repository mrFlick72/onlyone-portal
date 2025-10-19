from dataclasses import dataclass
import decimal
from decimal import ROUND_HALF_DOWN, Context, Decimal
from typing import Final

SCALE_PRECISION: Final = 2
SCALE_CRITERIA: Final = ROUND_HALF_DOWN

MONEY_ARITHMETIC_CONTEXT: Final = Context(rounding=SCALE_CRITERIA)

@dataclass
class Money:
    amount: decimal

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

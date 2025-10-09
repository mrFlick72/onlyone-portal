import decimal
from decimal import ROUND_HALF_DOWN, Context, Decimal
from typing import Final

SCALE_PRECISION: Final = 2
SCALE_CRITERIA: Final = ROUND_HALF_DOWN

MONEY_ARITHMETIC_CONTEXT: Final = Context(rounding=SCALE_CRITERIA)


class Money:

    def __init__(self, amount: decimal):
        self.amount = amount
        
    def __eq__(self, other):
        # Equality Comparison between two objects
        return self.amount == other.amount

    def __hash__(self):
        # hash(custom_object)
        return hash((self.amount))

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

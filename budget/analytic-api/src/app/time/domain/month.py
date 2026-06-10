from dataclasses import dataclass


@dataclass
class Month:
    content: int

    @staticmethod
    def JANUARY() -> "Month":
        return Month(1)

    @staticmethod
    def FEBRUARY() -> "Month":
        return Month(2)

    @staticmethod
    def MARCH() -> "Month":
        return Month(3)

    @staticmethod
    def APRIL() -> "Month":
        return Month(4)

    @staticmethod
    def MAY() -> "Month":
        return Month(5)

    @staticmethod
    def JUNE() -> "Month":
        return Month(6)

    @staticmethod
    def JULY() -> "Month":
        return Month(7)

    @staticmethod
    def AUGUST() -> "Month":
        return Month(8)

    @staticmethod
    def SEPTEMBER() -> "Month":
        return Month(9)

    @staticmethod
    def OCTOBER() -> "Month":
        return Month(10)

    @staticmethod
    def NOVEMBER() -> "Month":
        return Month(11)

    @staticmethod
    def DECEMBER() -> "Month":
        return Month(12)

from dataclasses import dataclass
from datetime import datetime, date
from typing import Final
from app.time.domain.year import Year
from app.time.domain.month import Month

DEFAULT_DATE_TIME_FORMATTER: Final = "%d/%m/%Y"
ISO_DATE_TIME_FORMATTER: Final = "%Y-%m-%d"


class DateParsingException(Exception):

    def __init__(self, message):
        self.message = message
        super().__init__(self.message)

@dataclass
class Date:
    content: date
    def formatted_date(self) -> str:
        return self.content.strftime(DEFAULT_DATE_TIME_FORMATTER)

    def iso_formatted_date(self) -> str:
        return self.content.strftime(ISO_DATE_TIME_FORMATTER)

    def plusYears(self, years: int):
        try:
            return Date(
                self.content.replace(year=self.content.year + years)
            )
        except ValueError:
            # Handle February 29 for leap years
            return Date(
                self.content.replace(month=2, day=28, year=self.content.year + years)
            )
    @staticmethod
    def date_for(content: str):
        try:
            return Date(datetime.strptime(content, DEFAULT_DATE_TIME_FORMATTER).date())
        except Exception as _:
            raise DateParsingException(
                f"wrong input formatting. It was {content} that it is not compatible with the pattern {DEFAULT_DATE_TIME_FORMATTER}"
            )

    @staticmethod
    def iso_date_for(content: str):
        try:
            return Date(datetime.strptime(content, ISO_DATE_TIME_FORMATTER).date())
        except Exception as _:
            raise DateParsingException(
                f"wrong input formatting. It was {content} that it is not compatible with the pattern {ISO_DATE_TIME_FORMATTER}"
            )

    @staticmethod
    def first_date_of_month(month: Month, year: Year):
        return Date.iso_date_for(f"{year.content}-{month.content}-01")

    @staticmethod
    def last_date_of_month(month: Month, year: Year):
        return Date.iso_date_for(
            f"{year.content}-{month.content}-{Date._last_day_for(month.content, year.content)}"
        )

    @staticmethod
    def _last_day_for(month: int, year: int):
        """returns the number of days in a given month"""
        days_per_month = [31, 28, 31, 30, 31, 30, 31, 31, 30, 31, 30, 31]
        d = days_per_month[month - 1]
        if month == 2 and (year % 4 == 0 and year % 100 != 0 or year % 400 == 0):
            d = 29
        return str(d)

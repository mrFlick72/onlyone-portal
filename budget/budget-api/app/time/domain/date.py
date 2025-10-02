from datetime import datetime, date

DEFAULT_DATE_TIME_FORMATTER = "%d/%m/%Y"
ISO_DATE_TIME_FORMATTER = "%Y-%m-%d"


class DateParsingException(Exception):

    def __init__(self, message):
        self.message = message
        super().__init__(self.message)


class Date:

    def __init__(self, content: date):
        self.content = content

    def formattedDate(self) -> str:
        return self.content.strftime(DEFAULT_DATE_TIME_FORMATTER)

    def isoFormattedDate(self) -> str:
        return self.content.strftime(ISO_DATE_TIME_FORMATTER)

    @staticmethod
    def dateFor(content: str):
        try:
            return Date(datetime.strptime(content, DEFAULT_DATE_TIME_FORMATTER).date())
        except Exception as _:
            raise DateParsingException(f"wrong input formatting. It was {content} that it is not compatible with the pattern {DEFAULT_DATE_TIME_FORMATTER}")

    @staticmethod
    def isoDateFor(content: str):
        try:
            return Date(datetime.strptime(content, ISO_DATE_TIME_FORMATTER).date())
        except Exception as _:
            raise DateParsingException(f"wrong input formatting. It was {content} that it is not compatible with the pattern {ISO_DATE_TIME_FORMATTER}")

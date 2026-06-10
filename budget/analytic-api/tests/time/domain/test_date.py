import pytest
from app.time.domain.date import Date, DateParsingException
from app.time.domain.month import Month
from app.time.domain.year import Year


def test_date_from_default_to_iso_format():
    expected_formatted_date = "2018-02-25"
    another_expected_formatted_fate = "2018-03-25"
    another_expected_formatted_date2 = "2018-05-25"

    date = Date.date_for("25/02/2018")
    anotherDate = Date.date_for("25/03/2018")
    anotherDate2 = Date.date_for("25/05/2018")

    assert date.iso_formatted_date() == expected_formatted_date
    assert anotherDate.iso_formatted_date() == another_expected_formatted_fate
    assert anotherDate2.iso_formatted_date() == another_expected_formatted_date2


def test_date_from_isp_to_default_format():
    expected_formatted_date = "25/02/2018"
    another_expected_formatted_fate = "25/03/2018"
    another_expected_formatted_date2 = "25/05/2018"

    date = Date.iso_date_for("2018-02-25")
    anotherDate = Date.iso_date_for("2018-03-25")
    anotherDate2 = Date.iso_date_for("2018-05-25")

    assert date.formatted_date() == expected_formatted_date
    assert anotherDate.formatted_date() == another_expected_formatted_fate
    assert anotherDate2.formatted_date() == another_expected_formatted_date2


def test_when_the_format_is_not_respected():
    with pytest.raises(DateParsingException):
        Date.iso_date_for("25/02/2018")


def test_first_date_of_month():
    expected = Date.date_for("01/02/2018")
    actual = Date.first_date_of_month(Month.FEBRUARY(), Year(2018))

    assert actual == expected


def test_last_date_of_month():
    expected = Date.date_for("28/02/2018")
    actual = Date.last_date_of_month(Month.FEBRUARY(), Year(2018))

    assert actual == expected

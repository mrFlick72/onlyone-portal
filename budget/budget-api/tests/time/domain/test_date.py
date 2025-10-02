"""package com.onlyoneportal.budget.time;

import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;

public class DateTest {
    static final DateTimeFormatter DATE_TIME_FORMATTER = DateTimeFormatter.ofPattern("dd-MM-YYY");

    @Test
    public void dateIsFormattedWithCustomFormatter() {
        String expectedFormattedDate = "25-02-2018";
        String anotherExpectedFormattedDate = "25-03-2018";
        String anotherExpectedFormattedDate2 = "25-05-2018";

        Date date = new Date(LocalDate.of(2018, 02, 25), DATE_TIME_FORMATTER);
        Date anotherDate = new Date(LocalDate.of(2018, 03, 25), DATE_TIME_FORMATTER);
        Date anotherDate2 = new Date(LocalDate.of(2018, 05, 25), DATE_TIME_FORMATTER);

        Assertions.assertEquals(date.formattedDate(), expectedFormattedDate);
        Assertions.assertEquals(anotherDate.formattedDate(), anotherExpectedFormattedDate);
        Assertions.assertEquals(anotherDate2.formattedDate(), anotherExpectedFormattedDate2);
    }


    @Test
    public void dateIsFormattedWithDefaultFormatter() {
        String expectedFormattedDate = "25/02/2018";
        String anotherExpectedFormattedDate = "25/03/2018";
        String anotherExpectedFormattedDate2 = "25/05/2018";

        Date date = new Date(LocalDate.of(2018, 02, 25));
        Date anotherDate = new Date(LocalDate.of(2018, 03, 25));
        Date anotherDate2 = new Date(LocalDate.of(2018, 05, 25));

        Assertions.assertEquals(date.formattedDate(), expectedFormattedDate);
        Assertions.assertEquals(anotherDate.formattedDate(), anotherExpectedFormattedDate);
        Assertions.assertEquals(anotherDate2.formattedDate(), anotherExpectedFormattedDate2);
    }

    @Test
    public void dateFromString() {
        Date expectedDateForDateString = new Date(LocalDate.of(2018, 02, 25));
        Date actualDateForDateString = Date.dateFor("25/02/2018");
        Assertions.assertEquals(actualDateForDateString, expectedDateForDateString);
    }

    @Test
    public void firstDateOfMonth() {
        Date expected = Date.dateFor("01/02/2018");
        Date actual = Date.firstDateOfMonth(Month.FEBRUARY, Year.of(2018));

        Assertions.assertEquals(actual, expected);
    }

    @Test
    public void lastDateOfMonth() {
        Date expected = Date.dateFor("28/02/2018");
        Date actual = Date.lastDateOfMonth(Month.FEBRUARY, Year.of(2018));

        Assertions.assertEquals(actual, expected);
    }
}"""

import pytest
from app.time.domain.date import Date, DateParsingException

date_time_formatter_pattern = "%Y-%m-%d"


def test_date_from_default_to_iso_format():
    expected_formatted_date = "2018-02-25"
    another_expected_formatted_fate = "2018-03-25"
    another_expected_formatted_date2 = "2018-05-25"

    date = Date.dateFor("25/02/2018")
    anotherDate = Date.dateFor("25/03/2018")
    anotherDate2 = Date.dateFor("25/05/2018")

    assert date.isoFormattedDate() == expected_formatted_date
    assert anotherDate.isoFormattedDate() == another_expected_formatted_fate
    assert anotherDate2.isoFormattedDate() == another_expected_formatted_date2

def test_date_from_isp_to_default_format():
    expected_formatted_date = "25/02/2018"
    another_expected_formatted_fate = "25/03/2018"
    another_expected_formatted_date2 = "25/05/2018"

    date = Date.isoDateFor("2018-02-25")
    anotherDate = Date.isoDateFor("2018-03-25")
    anotherDate2 = Date.isoDateFor("2018-05-25")

    assert date.formattedDate() == expected_formatted_date
    assert anotherDate.formattedDate() == another_expected_formatted_fate
    assert anotherDate2.formattedDate() == another_expected_formatted_date2

def test_when_the_format_is_not_respected():
    with pytest.raises(DateParsingException):
         Date.isoDateFor("25/02/2018")



def test_first_date_of_month():
    pass


def test_last_date_of_month():
    pass

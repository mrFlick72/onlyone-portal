package com.onlyoneportal.budget.time;

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
}
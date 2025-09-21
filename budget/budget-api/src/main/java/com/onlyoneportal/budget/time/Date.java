package com.onlyoneportal.budget.time;

import java.time.LocalDate;
import java.time.format.DateTimeFormatter;
import java.time.temporal.TemporalAdjusters;
import java.util.Objects;

import static java.time.format.DateTimeFormatter.ofPattern;


public final class Date implements Comparable<Date> {

    static final DateTimeFormatter DEFAULT_DATE_TIME_FORMATTER = ofPattern("dd/MM/yyyy");
    static final DateTimeFormatter ISO_DATE_TIME_FORMATTER = DateTimeFormatter.ISO_DATE;

    private final LocalDate localDate;
    private final DateTimeFormatter dateTimeFormatter;

    public Date(LocalDate localDate, DateTimeFormatter dateTimeFormatter) {
        this.localDate = localDate;
        this.dateTimeFormatter = dateTimeFormatter;
    }

    public Date(LocalDate localDate) {
        this.localDate = localDate;
        this.dateTimeFormatter = DEFAULT_DATE_TIME_FORMATTER;
    }

    public LocalDate getLocalDate() {
        return localDate;
    }

    public String formattedDate(){
        return dateTimeFormatter.format(localDate);
    }

    public String isoFormattedDate(){
        return ISO_DATE_TIME_FORMATTER.format(localDate);
    }

    public static Date dateFor(String date) {
        return new Date(LocalDate.parse(date, DEFAULT_DATE_TIME_FORMATTER));
    }
    public static Date isoDateFor(String date) {
        return new Date(LocalDate.parse(date, ISO_DATE_TIME_FORMATTER));
    }

    public static Date firstDateOfMonth(Month month, Year year) {
        return new Date(firstDateInLocalDate(month, year));
    }

    public static Date lastDateOfMonth(Month month, Year year) {
        return new Date(firstDateInLocalDate(month, year).with(TemporalAdjusters.lastDayOfMonth()));
    }

    private static LocalDate firstDateInLocalDate(Month month, Year year) {
        return LocalDate.of(year.getYearValue(), month.getMonthValue(), 1);
    }

    @Override
    public int compareTo(Date o) {
        return this.getLocalDate().compareTo(o.getLocalDate());
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Date date = (Date) o;
        return Objects.equals(localDate, date.localDate) && Objects.equals(dateTimeFormatter, date.dateTimeFormatter);
    }

    @Override
    public int hashCode() {
        return Objects.hash(localDate, dateTimeFormatter);
    }
}

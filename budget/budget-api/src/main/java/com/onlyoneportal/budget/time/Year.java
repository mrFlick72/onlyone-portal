package com.onlyoneportal.budget.time;

import java.util.Objects;

public final class Year {

    private final Integer yearValue;

    Year(Integer yearValue) {
        this.yearValue = yearValue;
    }

    public static Year of(Integer year) {
        return new Year(year);
    }

    public Integer getYearValue() {
        return yearValue;
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        Year year = (Year) o;
        return Objects.equals(yearValue, year.yearValue);
    }

    @Override
    public int hashCode() {
        return Objects.hash(yearValue);
    }
}

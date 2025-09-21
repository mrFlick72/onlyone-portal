package com.onlyoneportal.budget.expense.endpoint;



import java.io.Serializable;
import java.util.Objects;

public record TotalBySearchTagDetail(String searchTagKey, String searchTagValue, String total) implements Serializable {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        TotalBySearchTagDetail that = (TotalBySearchTagDetail) o;
        return Objects.equals(searchTagKey, that.searchTagKey) && Objects.equals(searchTagValue, that.searchTagValue) && Objects.equals(total, that.total);
    }

}

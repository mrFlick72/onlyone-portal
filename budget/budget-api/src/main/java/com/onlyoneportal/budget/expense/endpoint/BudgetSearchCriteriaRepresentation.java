package com.onlyoneportal.budget.expense.endpoint;

import java.io.Serializable;
import java.util.ArrayList;
import java.util.List;
import java.util.Objects;

public record BudgetSearchCriteriaRepresentation(Integer month, Integer year,
                                                 List<String> searchTagList) implements Serializable {

    public static BudgetSearchCriteriaRepresentation empty() {
        return new BudgetSearchCriteriaRepresentation(null, null, new ArrayList<>());
    }

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        BudgetSearchCriteriaRepresentation that = (BudgetSearchCriteriaRepresentation) o;
        return Objects.equals(month, that.month) && Objects.equals(year, that.year) && Objects.equals(searchTagList, that.searchTagList);
    }

}

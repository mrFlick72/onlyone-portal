package com.onlyoneportal.budget.expense.endpoint;


import java.util.List;
import java.util.Objects;

public record DailyBudgetExpenseRepresentation(List<BudgetExpenseRepresentation> budgetExpenseRepresentationList,
                                               String date, String total) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DailyBudgetExpenseRepresentation that = (DailyBudgetExpenseRepresentation) o;
        return Objects.equals(budgetExpenseRepresentationList, that.budgetExpenseRepresentationList) && Objects.equals(date, that.date) && Objects.equals(total, that.total);
    }

}

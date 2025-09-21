package com.onlyoneportal.budget.expense.model;

import java.util.Objects;

public record BudgetExpenseId(String content) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        BudgetExpenseId that = (BudgetExpenseId) o;
        return Objects.equals(content, that.content);
    }

}

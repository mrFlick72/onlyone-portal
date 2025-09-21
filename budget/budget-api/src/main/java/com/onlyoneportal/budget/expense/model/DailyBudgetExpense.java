package com.onlyoneportal.budget.expense.model;


import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.time.Date;

import java.util.List;
import java.util.Objects;

public record DailyBudgetExpense(List<BudgetExpense> budgetExpenseList, Date date, Money total) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        DailyBudgetExpense that = (DailyBudgetExpense) o;
        return Objects.equals(budgetExpenseList, that.budgetExpenseList) && Objects.equals(date, that.date) && Objects.equals(total, that.total);
    }

}

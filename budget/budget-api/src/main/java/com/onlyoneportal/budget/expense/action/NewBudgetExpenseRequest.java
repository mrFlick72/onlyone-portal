package com.onlyoneportal.budget.expense.action;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.time.Date;

import java.util.Objects;

public record NewBudgetExpenseRequest(Date date, Money amount, String note, String tag) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        NewBudgetExpenseRequest that = (NewBudgetExpenseRequest) o;
        return Objects.equals(date, that.date) && Objects.equals(amount, that.amount) && Objects.equals(note, that.note) && Objects.equals(tag, that.tag);
    }

}

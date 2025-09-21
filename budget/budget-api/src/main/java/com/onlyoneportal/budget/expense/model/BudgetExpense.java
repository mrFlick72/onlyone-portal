package com.onlyoneportal.budget.expense.model;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;

import java.util.Objects;

public record BudgetExpense(BudgetExpenseId id, UserName userName, Date date, Money amount, String note, String tag) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        BudgetExpense that = (BudgetExpense) o;
        return Objects.equals(id, that.id) && Objects.equals(userName, that.userName) && Objects.equals(date, that.date) && Objects.equals(amount, that.amount) && Objects.equals(note, that.note) && Objects.equals(tag, that.tag);
    }

}
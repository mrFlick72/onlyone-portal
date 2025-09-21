package com.onlyoneportal.budget.expense.endpoint;


import java.io.Serializable;
import java.util.Objects;

public record BudgetExpenseRepresentation(String id, String date, String amount, String note, String tagKey,
                                          String tagValue) implements Serializable {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        BudgetExpenseRepresentation that = (BudgetExpenseRepresentation) o;
        return Objects.equals(id, that.id) && Objects.equals(date, that.date) && Objects.equals(amount, that.amount) && Objects.equals(note, that.note) && Objects.equals(tagKey, that.tagKey) && Objects.equals(tagValue, that.tagValue);
    }

}
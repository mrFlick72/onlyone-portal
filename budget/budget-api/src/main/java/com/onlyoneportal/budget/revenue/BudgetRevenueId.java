package com.onlyoneportal.budget.revenue;

import java.util.Objects;

public record BudgetRevenueId(String content) {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        BudgetRevenueId that = (BudgetRevenueId) o;
        return Objects.equals(content, that.content);
    }

}

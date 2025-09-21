package com.onlyoneportal.budget.expense.endpoint;


import java.io.Serializable;
import java.util.List;
import java.util.Objects;

public record SpentBudgetRepresentation(String total,
                                        List<DailyBudgetExpenseRepresentation> dailyBudgetExpenseRepresentationList,
                                        List<TotalBySearchTagDetail> totalDetailList) implements Serializable {

    @Override
    public boolean equals(Object o) {
        if (this == o) return true;
        if (o == null || getClass() != o.getClass()) return false;
        SpentBudgetRepresentation that = (SpentBudgetRepresentation) o;
        return Objects.equals(dailyBudgetExpenseRepresentationList, that.dailyBudgetExpenseRepresentationList) && Objects.equals(totalDetailList, that.totalDetailList) && Objects.equals(total, that.total);
    }

}

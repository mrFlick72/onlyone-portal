package com.onlyoneportal.budget.expense.converter;

import com.onlyoneportal.budget.expense.endpoint.BudgetExpenseRepresentation;
import com.onlyoneportal.budget.expense.endpoint.DailyBudgetExpenseRepresentation;
import com.onlyoneportal.budget.expense.endpoint.SpentBudgetRepresentation;
import com.onlyoneportal.budget.expense.endpoint.TotalBySearchTagDetail;
import com.onlyoneportal.budget.expense.model.DailyBudgetExpense;
import com.onlyoneportal.budget.expense.model.SpentBudget;

import java.util.List;

import static java.util.stream.Collectors.toList;


public class SpentBudgetConverter {
    private final BudgetExpenseConverter budgetExpenseConverter;

    public SpentBudgetConverter(BudgetExpenseConverter budgetExpenseConverter) {
        this.budgetExpenseConverter = budgetExpenseConverter;
    }

    public SpentBudgetRepresentation domainToRepresentationModel(SpentBudget spentBudget) {
        return new SpentBudgetRepresentation(spentBudget.total().stringifyAmount(),
                spentBudget.dailyBudgetExpenseList().stream()
                        .map(dailyBudgetExpense ->
                                new DailyBudgetExpenseRepresentation(budgetExpenseRepresentationList(dailyBudgetExpense),
                                        dailyBudgetExpense.date().formattedDate(),
                                        dailyBudgetExpense.total().stringifyAmount())).collect(toList()),
                spentBudget.totalForSearchTags().entrySet().stream()
                        .map(total -> new TotalBySearchTagDetail(total.getKey().key(),
                                total.getKey().value(),
                                total.getValue().stringifyAmount()))
                        .collect(toList()));
    }

    private List<BudgetExpenseRepresentation> budgetExpenseRepresentationList(DailyBudgetExpense dailyBudgetExpense) {
        return dailyBudgetExpense.budgetExpenseList().stream()
                .map(budgetExpenseConverter::domainToRepresentationModel)
                .collect(toList());
    }


}

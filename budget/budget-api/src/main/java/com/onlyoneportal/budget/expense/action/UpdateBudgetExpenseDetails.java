package com.onlyoneportal.budget.expense.action;

import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;

public class UpdateBudgetExpenseDetails {
    private final BudgetExpenseRepository budgetExpenseRepository;

    public UpdateBudgetExpenseDetails(BudgetExpenseRepository budgetExpenseRepository) {

        this.budgetExpenseRepository = budgetExpenseRepository;
    }

    public void updateWithoutAttachment(BudgetExpense budgetExpense) {
        budgetExpenseRepository.findFor(budgetExpense.id())
                .ifPresent(foundBudgetExpense -> {
                    BudgetExpense updatedBudgetExpense = new BudgetExpense(budgetExpense.id(),
                            budgetExpense.userName(),
                            budgetExpense.date(),
                            budgetExpense.amount(), budgetExpense.note(),
                            budgetExpense.tag()
                    );

                    budgetExpenseRepository.save(updatedBudgetExpense);
                });
    }
}

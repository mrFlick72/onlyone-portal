package com.onlyoneportal.budget.expense.action;

import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.expense.model.BudgetExpenseNotFoundException;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;

public class DeleteBudgetExpense {

    private final BudgetExpenseRepository budgetExpenseRepository;

    public DeleteBudgetExpense(BudgetExpenseRepository budgetExpenseRepository) {

        this.budgetExpenseRepository = budgetExpenseRepository;
    }

    public void delete(BudgetExpenseId budgetExpenseId) {
        budgetExpenseRepository.findFor(budgetExpenseId)
                .ifPresentOrElse(budgetExpense -> budgetExpenseRepository.delete(budgetExpense.id()),
                        () -> {
                            throw new BudgetExpenseNotFoundException();
                        });
    }
}

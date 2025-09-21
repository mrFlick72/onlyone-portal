package com.onlyoneportal.budget.expense.action;

import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.user.UserRepository;

import java.util.Optional;

public class CreateBudgetExpense {

    private final UserRepository userRepository;
    private final BudgetExpenseRepository budgetExpenseRepository;
    private final String DEFAULT_TAG = SearchTag.DEFAULT_KEY;

    public CreateBudgetExpense(BudgetExpenseRepository budgetExpenseRepository,
                               UserRepository userRepository) {
        this.budgetExpenseRepository = budgetExpenseRepository;
        this.userRepository = userRepository;
    }

    public BudgetExpense newBudgetExpense(NewBudgetExpenseRequest newBudgetExpenseRequest) {
        BudgetExpense budgetExpense = new BudgetExpense(null,
                userRepository.currentLoggedUserName(),
                newBudgetExpenseRequest.date(),
                newBudgetExpenseRequest.amount(),
                newBudgetExpenseRequest.note(),
                getSearchTag(newBudgetExpenseRequest));

        budgetExpenseRepository.save(budgetExpense);
        return budgetExpense;
    }

    private String getSearchTag(NewBudgetExpenseRequest newBudgetExpenseRequest) {
        return Optional.ofNullable(newBudgetExpenseRequest.tag())
                .map(tag -> tag.isBlank() ? DEFAULT_TAG : tag)
                .orElse(DEFAULT_TAG);
    }
}

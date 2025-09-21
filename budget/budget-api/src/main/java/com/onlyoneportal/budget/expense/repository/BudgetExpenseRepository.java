package com.onlyoneportal.budget.expense.repository;

import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;

import java.util.List;
import java.util.Optional;

public interface BudgetExpenseRepository {

    Optional<BudgetExpense> findFor(BudgetExpenseId budgetExpenseId);

    List<BudgetExpense> findByDateRange(UserName userName, Date star, Date end, String... searchTags);

    BudgetExpense save(BudgetExpense budgetExpense);

    void delete(BudgetExpenseId idBudgetExpense);
}

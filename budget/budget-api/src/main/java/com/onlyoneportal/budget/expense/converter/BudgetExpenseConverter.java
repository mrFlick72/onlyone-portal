package com.onlyoneportal.budget.expense.converter;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.action.NewBudgetExpenseRequest;
import com.onlyoneportal.budget.expense.endpoint.BudgetExpenseRepresentation;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.searchtag.SearchTagRepository;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserRepository;

import java.util.Optional;

public class BudgetExpenseConverter {

    private final SearchTagRepository searchTagRepository;
    private final UserRepository userRepository;

    public BudgetExpenseConverter(SearchTagRepository searchTagRepository, UserRepository userRepository) {
        this.searchTagRepository = searchTagRepository;
        this.userRepository = userRepository;
    }

    public BudgetExpenseRepresentation domainToRepresentationModel(BudgetExpense budgetExpense) {
        String searchTag = Optional.ofNullable(searchTagRepository.findSearchTagBy(budgetExpense.tag()))
                .map(SearchTag::value).orElse("");
        return new BudgetExpenseRepresentation(budgetExpense.id().content(), budgetExpense.date().formattedDate(),
                budgetExpense.amount().stringifyAmount(), budgetExpense.note(), budgetExpense.tag(), searchTag);
    }

    public BudgetExpense representationModelToDomainModel(BudgetExpenseRepresentation budgetExpenseRepresentation) {
        return new BudgetExpense(new BudgetExpenseId(budgetExpenseRepresentation.id()),
                userRepository.currentLoggedUserName(),
                Date.dateFor(budgetExpenseRepresentation.date()),
                Money.moneyFor(budgetExpenseRepresentation.amount()),
                budgetExpenseRepresentation.note(), budgetExpenseRepresentation.tagKey());
    }

    public NewBudgetExpenseRequest newBudgetExpenseRequestFromRepresentation(BudgetExpenseRepresentation budgetExpenseRepresentation) {
        return new NewBudgetExpenseRequest(Date.dateFor(budgetExpenseRepresentation.date()),
                Money.moneyFor(budgetExpenseRepresentation.amount()),
                budgetExpenseRepresentation.note(), budgetExpenseRepresentation.tagKey());
    }

}

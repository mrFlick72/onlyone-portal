package com.onlyoneportal.budget.expense.action;

import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.SpentBudget;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.searchtag.SearchTagRepository;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.time.Month;
import com.onlyoneportal.budget.time.Year;
import com.onlyoneportal.budget.user.UserName;
import com.onlyoneportal.budget.user.UserRepository;

import java.util.Comparator;
import java.util.List;
import java.util.Optional;

import static java.util.stream.Collectors.toList;

public class FindSpentBudget {

    private final UserRepository userRepository;
    private final BudgetExpenseRepository budgetExpenseRepository;
    private final SearchTagRepository searchTagRepository;

    public FindSpentBudget(UserRepository userRepository,
                           BudgetExpenseRepository budgetExpenseRepository,
                           SearchTagRepository searchTagRepository) {
        this.userRepository = userRepository;
        this.budgetExpenseRepository = budgetExpenseRepository;
        this.searchTagRepository = searchTagRepository;
    }

    public SpentBudget findBy(Month month, Year year, List<String> searchTagList) {
        UserName userName = userRepository.currentLoggedUserName();
        String[] searchTags = searchTags(searchTagList);

        List<BudgetExpense> budgetExpenseList =
                budgetExpenseRepository.findByDateRange(userName, Date.firstDateOfMonth(month, year),
                        Date.lastDateOfMonth(month, year),
                        searchTags);

        List<SearchTag> allSearchTagFor = getAllSearchTagFor(budgetExpenseList);

        return new SpentBudget(orderByDate(budgetExpenseList), allSearchTagFor);
    }

    private String[] searchTags(List<String> searchTagList) {
        return Optional.ofNullable(searchTagList)
                .map(searchTagListAux -> searchTagListAux.stream()
                        .filter(tag -> !"".equals(tag.trim()))
                        .collect(toList()))
                .map(tags -> tags.toArray(new String[tags.size()]))
                .orElse(new String[0]);
    }

    private List<BudgetExpense> orderByDate(List<BudgetExpense> budgetExpenseListFromRepository) {
        return budgetExpenseListFromRepository.stream()
                .sorted(Comparator.comparing(BudgetExpense::date))
                .collect(toList());
    }

    private List<SearchTag> getAllSearchTagFor(List<BudgetExpense> budgetExpenses) {
        return budgetExpenses.stream().map(BudgetExpense::tag)
                .distinct().map(searchTagRepository::findSearchTagBy)
                .collect(toList());
    }
}
package com.onlyoneportal.budget.config;

import com.onlyoneportal.budget.expense.action.CreateBudgetExpense;
import com.onlyoneportal.budget.expense.action.DeleteBudgetExpense;
import com.onlyoneportal.budget.expense.action.FindSpentBudget;
import com.onlyoneportal.budget.expense.action.UpdateBudgetExpenseDetails;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;
import com.onlyoneportal.budget.revenue.BudgetRevenueRepository;
import com.onlyoneportal.budget.revenue.FindBudgetRevenue;
import com.onlyoneportal.budget.searchtag.SearchTagRepository;
import com.onlyoneportal.budget.user.UserRepository;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration(proxyBeanMethods = false)
public class BusinessLogicConfiguration {

    @Bean
    public DeleteBudgetExpense deleteBudgetExpense(BudgetExpenseRepository budgetExpenseRepository) {
        return new DeleteBudgetExpense(budgetExpenseRepository);
    }

    @Bean
    public UpdateBudgetExpenseDetails updateBudgetExpenseDetails(BudgetExpenseRepository budgetExpenseRepository) {
        return new UpdateBudgetExpenseDetails(budgetExpenseRepository);
    }


    @Bean
    public FindBudgetRevenue findBudgetRevenue(UserRepository userRepository,
                                               BudgetRevenueRepository budgetRevenueRepository) {
        return new FindBudgetRevenue(budgetRevenueRepository, userRepository);
    }

    @Bean
    public FindSpentBudget findSpentBudget(UserRepository userRepository,
                                           SearchTagRepository searchTagRepository,
                                           BudgetExpenseRepository budgetExpenseRepository) {
        return new FindSpentBudget(userRepository, budgetExpenseRepository, searchTagRepository);
    }


    @Bean
    public CreateBudgetExpense createBudgetExpense(UserRepository userRepository,
                                                   BudgetExpenseRepository budgetExpenseRepository) {
        return new CreateBudgetExpense(budgetExpenseRepository, userRepository);
    }

}

package com.onlyoneportal.budget.config;

import com.onlyoneportal.budget.expense.converter.BudgetExpenseConverter;
import com.onlyoneportal.budget.expense.converter.SpentBudgetConverter;
import com.onlyoneportal.budget.revenue.BudgetRevenueConverter;
import com.onlyoneportal.budget.searchtag.SearchTagRepository;
import com.onlyoneportal.budget.user.UserRepository;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

@Configuration(proxyBeanMethods = false)
public class AdapterConfiguration {

    @Bean
    public BudgetRevenueConverter budgetRevenueAdapter(UserRepository userRepository) {
        return new BudgetRevenueConverter(userRepository);
    }

    @Bean
    public BudgetExpenseConverter budgetExpenseAdapter(UserRepository userRepository,
                                                       SearchTagRepository searchTagRepository) {
        return new BudgetExpenseConverter(searchTagRepository, userRepository);
    }

    @Bean
    public SpentBudgetConverter spentBudgetAdapter(BudgetExpenseConverter budgetExpenseConverter) {
        return new SpentBudgetConverter(budgetExpenseConverter);
    }
}
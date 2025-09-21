package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserRepository;


public class BudgetRevenueConverter {

    private final UserRepository userRepository;

    public BudgetRevenueConverter(UserRepository userRepository) {
        this.userRepository = userRepository;
    }

    public BudgetRevenue fromRepresentationToModel(BudgetRevenueRepresentation budgetRevenueRepresentation) {
        return new BudgetRevenue(
                new BudgetRevenueId(budgetRevenueRepresentation.id()),
                userRepository.currentLoggedUserName().content(),
                Date.dateFor(budgetRevenueRepresentation.date()),
                Money.moneyFor(budgetRevenueRepresentation.amount()),
                budgetRevenueRepresentation.note());
    }

    public BudgetRevenueRepresentation fromDomainToRepresentation(BudgetRevenue budgetRevenue) {
        return new BudgetRevenueRepresentation(budgetRevenue.id().content(),
                budgetRevenue.registrationDate().formattedDate(),
                budgetRevenue.amount().stringifyAmount(),
                budgetRevenue.note());
    }
}

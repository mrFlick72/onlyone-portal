package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.time.Month;
import com.onlyoneportal.budget.time.Year;
import com.onlyoneportal.budget.user.UserRepository;

import java.util.List;

public class FindBudgetRevenue {

    private final BudgetRevenueRepository budgetRevenueRepository;
    private final UserRepository userRepository;

    public FindBudgetRevenue(BudgetRevenueRepository budgetRevenueRepository, UserRepository userRepository) {
        this.budgetRevenueRepository = budgetRevenueRepository;
        this.userRepository = userRepository;
    }

    public List<BudgetRevenue> findBy(Year year) {
        return budgetRevenueRepository.findByDateRange(userRepository.currentLoggedUserName().content(),
                Date.firstDateOfMonth(Month.JANUARY, year),
                Date.lastDateOfMonth(Month.DECEMBER, year));
    }
}
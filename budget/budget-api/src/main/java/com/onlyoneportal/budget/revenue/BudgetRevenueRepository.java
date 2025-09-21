package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.time.Date;

import java.util.List;

public interface BudgetRevenueRepository {

    List<BudgetRevenue> findByDateRange(String user, Date star, Date end);

    BudgetRevenue save(BudgetRevenue budgetRevenue);

    void update(BudgetRevenue budgetRevenue);

    void delete(BudgetRevenueId idBudgetRevenue);
}

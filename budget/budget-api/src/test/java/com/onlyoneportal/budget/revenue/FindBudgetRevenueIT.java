package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.time.Year;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.springframework.beans.factory.annotation.Autowired;
import org.springframework.boot.test.context.SpringBootTest;
import org.springframework.security.test.context.support.WithMockUser;
import org.springframework.test.context.jdbc.Sql;
import org.springframework.test.context.jdbc.SqlGroup;
import org.springframework.test.context.junit.jupiter.SpringExtension;

import java.util.List;

import static java.util.Arrays.asList;


@SpringBootTest
@ExtendWith(SpringExtension.class)
public class FindBudgetRevenueIT {

    @Autowired
    private FindBudgetRevenue findBudgetRevenue;

    @Test
    @SqlGroup(value = {
            @Sql("classpath:schema.sql"),
            @Sql("classpath:budget_revenue/find-year-data-set.sql")
    })
    @WithMockUser("USER")
    public void findYearlyBudgetRevenue() {
        Year of = Year.of(2018);

        List<BudgetRevenue> yearlyBudgetRevenueList = findBudgetRevenue.findBy(of);
        List<BudgetRevenue> expectedBudgetRevenueList =
                asList(new BudgetRevenue(null, "USER", Date.dateFor("01/01/2018"), Money.ONE, ""),
                        new BudgetRevenue(null, "USER", Date.dateFor("01/02/2018"), Money.ONE, ""),
                        new BudgetRevenue(null, "USER", Date.dateFor("01/03/2018"), Money.ONE, ""),
                        new BudgetRevenue(null, "USER", Date.dateFor("01/05/2018"), Money.ONE, ""),
                        new BudgetRevenue(null, "USER", Date.dateFor("01/06/2018"), Money.ONE, ""),
                        new BudgetRevenue(null, "USER", Date.dateFor("01/09/2018"), Money.ONE, ""),
                        new BudgetRevenue(null, "USER", Date.dateFor("01/10/2018"), Money.ONE, ""));

        System.out.println(yearlyBudgetRevenueList);
        System.out.println(expectedBudgetRevenueList);
        Assertions.assertEquals(yearlyBudgetRevenueList, expectedBudgetRevenueList);
    }
}
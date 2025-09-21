package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.DailyBudgetExpense;
import com.onlyoneportal.budget.expense.model.SpentBudget;
import com.onlyoneportal.budget.support.BudgetFixture;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import java.util.List;

import static java.util.Arrays.asList;

public class SpentBudgetRevenueTest {

    @Test
    public void getTotal() {
        SpentBudget spent25Budget =
                new SpentBudget(asList(new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market")), asList());
        SpentBudget spent30Budget =
                new SpentBudget(asList(new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("5"), "super-market", "super-market")), asList());

        SpentBudget spent50Budget =
                new SpentBudget(asList(new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("20.00"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("5"), "super-market", "super-market")), asList());

        Assertions.assertEquals(Money.moneyFor("30.00"), spent30Budget.total());
        Assertions.assertEquals(Money.moneyFor("25.00"), spent25Budget.total());
        Assertions.assertEquals(Money.moneyFor("50.00"), spent50Budget.total());
    }

    @Test
    public void dailyBudgetExpenseList() {
        SpentBudget spentBudget =
                new SpentBudget(asList(new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("15/02/2018"), Money.moneyFor("20.00"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("15/02/2018"), Money.moneyFor("5"), "super-market", "super-market")), asList());

        List<DailyBudgetExpense> expected = asList(
                new DailyBudgetExpense(asList(new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market")),
                        Date.dateFor("12/02/2018"), Money.moneyFor("25.00")),
                new DailyBudgetExpense(asList(new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("15/02/2018"), Money.moneyFor("20.00"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), new UserName("USER"), Date.dateFor("15/02/2018"), Money.moneyFor("5"), "super-market", "super-market")),
                        Date.dateFor("15/02/2018"), Money.moneyFor("25.00")));

        Assertions.assertEquals(spentBudget.dailyBudgetExpenseList(), expected);
    }

}
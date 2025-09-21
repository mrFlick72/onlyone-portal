package com.onlyoneportal.budget.infrastructure.dynamodb;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;

import static com.onlyoneportal.budget.support.BudgetFixture.saltGenerator;

class BudgetDynamoDbIdFactoryTest {

    public static final BudgetExpense BUDGET_EXPENSE = new BudgetExpense(
            null,
            new UserName("USER"),
            Date.dateFor("01/01/2018"),
            Money.ONE,
            "",
            ""
    );


    @Test
    void getACompleteBudgetId() {
        BudgetDynamoDbIdFactory<BudgetExpenseId, BudgetExpense> budgetDynamoDbIdFactory = new BudgetExpenseDynamoDbIdFactory(saltGenerator);
        BudgetExpenseId actual = budgetDynamoDbIdFactory.budgetIdFrom(BUDGET_EXPENSE);

        Assertions.assertEquals("MjAxOF8xX1VTRVI=-MV9BX1NBTFQ=", actual.content());
    }
}
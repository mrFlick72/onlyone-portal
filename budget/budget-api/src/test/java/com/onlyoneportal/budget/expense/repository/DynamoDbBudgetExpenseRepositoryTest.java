package com.onlyoneportal.budget.expense.repository;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetExpenseDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import com.onlyoneportal.budget.user.UserRepository;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.List;
import java.util.Optional;

import static com.onlyoneportal.budget.support.BudgetFixture.saltGenerator;
import static com.onlyoneportal.budget.support.DatabaseUtils.*;
import static com.onlyoneportal.budget.support.UserTestFixture.A_USER_NAME;
import static java.util.Arrays.asList;
import static org.mockito.BDDMockito.given;

@ExtendWith(MockitoExtension.class)
class DynamoDbBudgetExpenseRepositoryTest {
    private static final Date DATE = Date.dateFor("12/02/2018");

    @Mock
    private UserRepository userRepository;

    private BudgetExpenseRepository budgetExpenseRepository;
    private final BudgetDynamoDbIdFactory idFactory = new BudgetExpenseDynamoDbIdFactory(saltGenerator);

    private void databaseInit() {
        loadBudgetExpense().forEach(budgetExpenseRepository::save);
    }

    @BeforeEach
    public void setUp() {
        budgetExpenseRepository = new DynamoDbBudgetExpenseRepository(BUDGET_EXPENSE_TABLE_NAME,
                dynamoClient(),
                idFactory,
                userRepository,
                new DynamoDbAttributeValueFactory()
        );
        given(userRepository.currentLoggedUserName())
                .willReturn(A_USER_NAME);

        resetDatabase();
        databaseInit();
    }

    @Test
    public void saveAnewBudgetExpense() {
        BudgetExpenseId id = new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU");
        BudgetExpense expected = new BudgetExpense(id, new UserName("USER"), DATE, Money.moneyFor("10.50"), "NOTE", "TAG");

        BudgetExpense actual = budgetExpenseRepository.save(expected);

        Assertions.assertEquals(expected, actual);
        BudgetExpense retrievedBudgetExpense = budgetExpenseRepository.findFor(expected.id()).get();
        Assertions.assertEquals(expected, retrievedBudgetExpense);
    }

    @Test
    public void findByDateRange() {
        budgetExpenseRepository.save(new BudgetExpense(null, new UserName("USER"), Date.dateFor("06/05/2018"), Money.moneyFor("17.50"), "Lanch", "lanch"));
        List<BudgetExpense> actualRange = budgetExpenseRepository.findByDateRange(new UserName("USER"), Date.dateFor("01/01/2018"), Date.dateFor("05/05/2018"));
        List<BudgetExpense> expectedRange =
                asList(
                        new BudgetExpense(null, new UserName("USER"), Date.dateFor("06/01/2018"), Money.moneyFor("17.50"), "Lanch", "lanch"),
                        new BudgetExpense(null, new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("10.50"), "Super Market", "super-market"),
                        new BudgetExpense(null, new UserName("USER"), Date.dateFor("13/02/2018"), Money.moneyFor("17.50"), "Dinner", "dinner"),
                        new BudgetExpense(null, new UserName("USER"), Date.dateFor("22/02/2018"), Money.moneyFor("17.50"), "Super Market", "super-market"),
                        new BudgetExpense(null, new UserName("USER"), Date.dateFor("05/05/2018"), Money.moneyFor("17.50"), "Lanch", "lanch")
                );

        Assertions.assertEquals(expectedRange.size(), actualRange.size());
    }

    @Test
    public void findByDateRangeAndSearchTags() {
        List<BudgetExpense> actualRange = budgetExpenseRepository.findByDateRange(new UserName("USER"), Date.dateFor("01/01/2018"), Date.dateFor("05/05/2018"), "super-market", "dinner");
        List<BudgetExpense> expectedRange =
                asList(new BudgetExpense(new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU"), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("10.50"), "Super Market", "super-market"),
                        new BudgetExpense(new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTNfQV9TQUxU"), new UserName("USER"), Date.dateFor("13/02/2018"), Money.moneyFor("17.50"), "Dinner", "dinner"),
                        new BudgetExpense(new BudgetExpenseId("MjAxOF8yX1VTRVI=-MjJfQV9TQUxU"), new UserName("USER"), Date.dateFor("22/02/2018"), Money.moneyFor("17.50"), "Super Market", "super-market"));

        Assertions.assertTrue(expectedRange.containsAll(actualRange) );
    }
    @Test
    public void deleteBudgetExpense() {
        BudgetExpenseId id = new BudgetExpenseId("MjAxOF8yX1VTRVI=-MTJfQV9TQUxU");
        BudgetExpense expected = new BudgetExpense(id, new UserName("USER"), DATE, Money.moneyFor("10.50"), "NOTE", "TAG");

        budgetExpenseRepository.save(expected);

        budgetExpenseRepository.delete(id);
        Optional<BudgetExpense> actual = budgetExpenseRepository.findFor(id);

        Assertions.assertThrows(Exception.class, actual::orElseThrow);
    }

}
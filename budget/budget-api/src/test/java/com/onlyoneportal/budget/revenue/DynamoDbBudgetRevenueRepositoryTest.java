package com.onlyoneportal.budget.revenue;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.BudgetRevenueDynamoDbIdFactory;
import com.onlyoneportal.budget.infrastructure.dynamodb.DynamoDbAttributeValueFactory;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserRepository;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.BeforeEach;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;
import software.amazon.awssdk.services.dynamodb.DynamoDbClient;

import java.util.List;

import static com.onlyoneportal.budget.support.BudgetFixture.saltGenerator;
import static com.onlyoneportal.budget.support.DatabaseUtils.*;
import static com.onlyoneportal.budget.support.UserTestFixture.A_USER_NAME;
import static java.util.Arrays.asList;
import static org.mockito.BDDMockito.given;

@ExtendWith(MockitoExtension.class)
class DynamoDbBudgetRevenueRepositoryTest {

    private static final Date DATE = Date.dateFor("12/02/2000");
    public static final DynamoDbClient DYNAMO_DB_CLIENT = dynamoClient();

    @Mock
    private UserRepository userRepository;

    private BudgetRevenueRepository budgetRevenueRepository;

    private final BudgetDynamoDbIdFactory idFactory = new BudgetRevenueDynamoDbIdFactory(saltGenerator);


    @BeforeEach
    public void setUp() {
        budgetRevenueRepository = new DynamoDbBudgetRevenueRepository(BUDGET_REVENUE_TABLE_NAME,
                DYNAMO_DB_CLIENT,
                idFactory,
                userRepository,
                new DynamoDbAttributeValueFactory());

        given(userRepository.currentLoggedUserName())
                .willReturn(A_USER_NAME);

        resetDatabase();
    }

    private void databaseInit() {
        loadBudgetRevenue().forEach(budgetRevenueRepository::save);
    }

    @Test
    void saveAnewBudgetRevenue() {
        BudgetRevenue budgetRevenue = new BudgetRevenue(null, "USER", DATE, Money.ONE, "A_NOTE");

        budgetRevenueRepository.save(budgetRevenue);

        List<BudgetRevenue> expectedRange = budgetRevenueRepository.findByDateRange("USER", DATE, DATE);

        Assertions.assertFalse(expectedRange.isEmpty());
    }

    @Test
    public void findByDateRange() {
        databaseInit();

        List<BudgetRevenue> actualRange = budgetRevenueRepository.findByDateRange("USER", Date.dateFor("01/01/2018"), Date.dateFor("05/05/2018"));
        List<BudgetRevenue> expectedRange =
                asList(new BudgetRevenue(new BudgetRevenueId("MjAxOF9VU0VS-MV82X0FfU0FMVA=="), "USER", Date.dateFor("06/01/2018"), Money.moneyFor("17.50"), "Lanch"),
                        new BudgetRevenue(new BudgetRevenueId("MjAxOF9VU0VS-Ml8xMl9BX1NBTFQ="), "USER", Date.dateFor("12/02/2018"), Money.moneyFor("10.50"), "Super Market"),
                        new BudgetRevenue(new BudgetRevenueId("MjAxOF9VU0VS-Ml8xM19BX1NBTFQ="), "USER", Date.dateFor("13/02/2018"), Money.moneyFor("17.50"), "Dinner"),
                        new BudgetRevenue(new BudgetRevenueId("MjAxOF9VU0VS-Ml8yMl9BX1NBTFQ="), "USER", Date.dateFor("22/02/2018"), Money.moneyFor("17.50"), "Super Market"));

        System.out.println(actualRange);
        Assertions.assertTrue(actualRange.containsAll(expectedRange));
    }


    @Test
    public void updateBudgetRevenue() {
        BudgetRevenue budgetRevenue = new BudgetRevenue(null, "USER", DATE, Money.ONE, "A_NOTE");
        BudgetRevenue revenue = budgetRevenueRepository.save(budgetRevenue);


        budgetRevenueRepository.update(new BudgetRevenue(
                revenue.id(), revenue.userName(),
                revenue.registrationDate(),
                revenue.amount(),
                "A NEW NOTE"
        ));

        List<BudgetRevenue> user = budgetRevenueRepository.findByDateRange("USER", DATE, DATE);

        Assertions.assertEquals("A NEW NOTE", user.get(0).note());
    }

    @Test
    public void deleteBudgetExpense() {
        BudgetRevenue budgetRevenue = new BudgetRevenue(null, "USER", DATE, Money.ONE, "A_NOTE");

        BudgetRevenue revenue = budgetRevenueRepository.save(budgetRevenue);

        List<BudgetRevenue> expectedRange = budgetRevenueRepository.findByDateRange("USER", DATE, DATE);

        Assertions.assertFalse(expectedRange.isEmpty());

        budgetRevenueRepository.delete(revenue.id());
        expectedRange = budgetRevenueRepository.findByDateRange("USER", DATE, DATE);

        Assertions.assertTrue(expectedRange.isEmpty());
    }

}
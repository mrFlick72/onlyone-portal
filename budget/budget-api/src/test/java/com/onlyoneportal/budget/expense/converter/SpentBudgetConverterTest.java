package com.onlyoneportal.budget.expense.converter;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.endpoint.BudgetExpenseRepresentation;
import com.onlyoneportal.budget.expense.endpoint.DailyBudgetExpenseRepresentation;
import com.onlyoneportal.budget.expense.endpoint.SpentBudgetRepresentation;
import com.onlyoneportal.budget.expense.endpoint.TotalBySearchTagDetail;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.expense.model.SpentBudget;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.searchtag.SearchTagRepository;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static java.util.Arrays.asList;
import static org.mockito.BDDMockito.given;

@ExtendWith(MockitoExtension.class)
public class SpentBudgetConverterTest {

    private static final String DATE_12_02_2018 = "12/02/2018";
    private static final String DATE_22_02_2018 = "22/02/2018";

    @Mock
    private SearchTagRepository searchTagRepository;

    @Test
    public void spentBudgetIsConverted() {
        SpentBudgetRepresentation expected = getExpectation();
        given(searchTagRepository.findSearchTagBy("super-market")).willReturn(new SearchTag("super-market", "super-market"));

        SpentBudgetRepresentation actual = new SpentBudgetConverter(new BudgetExpenseConverter(searchTagRepository, null))
                .domainToRepresentationModel(new SpentBudget(asList(
                        new BudgetExpense(new BudgetExpenseId("1"), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("12.50"), "", "super-market"),
                        new BudgetExpense(new BudgetExpenseId("2"), new UserName("USER"), Date.dateFor("12/02/2018"), Money.moneyFor("15.50"), "", "super-market"),
                        new BudgetExpense(new BudgetExpenseId("3"), new UserName("USER"), Date.dateFor("22/02/2018"), Money.moneyFor("22.50"), "", "super-market")),
                        asList(new SearchTag("super-market", "super-market"))));


        Assertions.assertEquals(actual, expected);
    }

    private SpentBudgetRepresentation getExpectation() {
        return new SpentBudgetRepresentation("50.50", asList(new DailyBudgetExpenseRepresentation(asList(new BudgetExpenseRepresentation("1", DATE_12_02_2018, "12.50", "", "super-market", "super-market"),
                        new BudgetExpenseRepresentation("2", DATE_12_02_2018, "15.50", "", "super-market", "super-market")), DATE_12_02_2018, "28.00"),
                new DailyBudgetExpenseRepresentation(asList(new BudgetExpenseRepresentation("3", DATE_22_02_2018, "22.50", "", "super-market", "super-market")), DATE_22_02_2018, "22.50")),
                asList(new TotalBySearchTagDetail("super-market", "super-market", "50.50")));
    }
}
package com.onlyoneportal.budget.expense.converter;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.endpoint.BudgetExpenseRepresentation;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.searchtag.SearchTagRepository;
import com.onlyoneportal.budget.support.BudgetFixture;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import com.onlyoneportal.budget.user.UserRepository;
import org.junit.jupiter.api.Assertions;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.mockito.BDDMockito.given;
import static org.mockito.Mockito.verify;

@ExtendWith(MockitoExtension.class)
public class BudgetExpenseConverterTest {

    private static final String AMOUNT = "12.50";
    private static final Money MONEY_AMOUNT = Money.moneyFor("12.50");
    private static final String DATE = "25/02/2018";
    private static final Date DOMAIN_DATE = Date.dateFor("25/02/2018");
    private final UserName USER = new UserName("USER");

    @Mock
    private SearchTagRepository searchTagRepository;

    @Mock
    private UserRepository userRepository;

    @Test
    public void convertWebRepresentationToDomainModel() {
        BudgetExpenseId id = BudgetFixture.randomBudgetExpenseId();
        BudgetExpenseRepresentation budgetExpenseRepresentation = new BudgetExpenseRepresentation(id.content(), DATE, AMOUNT, "super-market", "super-market", "Super Market");
        BudgetExpenseConverter budgetExpenseConverter = new BudgetExpenseConverter(searchTagRepository, userRepository);

        given(userRepository.currentLoggedUserName())
                .willReturn(USER);

        BudgetExpense actualBudgetExpense = budgetExpenseConverter.representationModelToDomainModel(budgetExpenseRepresentation);
        verify(userRepository).currentLoggedUserName();

        Assertions.assertEquals(actualBudgetExpense, new BudgetExpense(id, USER, DOMAIN_DATE, MONEY_AMOUNT, "super-market", "super-market"));
    }

    @Test
    public void convertDomainModelToWebRepresentation() {
        BudgetExpenseId id = BudgetFixture.randomBudgetExpenseId();

        BudgetExpense budgetExpense = new BudgetExpense(id, USER, DOMAIN_DATE, MONEY_AMOUNT, "Super Market", "super-market");
        BudgetExpenseConverter budgetExpenseConverter = new BudgetExpenseConverter(searchTagRepository, userRepository);

        given(searchTagRepository.findSearchTagBy("super-market"))
                .willReturn(new SearchTag("super-market", "Super Market"));

        BudgetExpenseRepresentation actualBudgetExpenseRepresentation =
                budgetExpenseConverter.domainToRepresentationModel(budgetExpense);

        verify(searchTagRepository).findSearchTagBy("super-market");

        Assertions.assertEquals(actualBudgetExpenseRepresentation, new BudgetExpenseRepresentation(id.content(), DATE, AMOUNT, "Super Market", "super-market", "Super Market"));
    }

}
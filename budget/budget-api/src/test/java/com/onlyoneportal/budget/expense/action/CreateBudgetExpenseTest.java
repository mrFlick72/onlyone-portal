package com.onlyoneportal.budget.expense.action;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import com.onlyoneportal.budget.user.UserRepository;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import static org.mockito.BDDMockito.given;
import static org.mockito.Mockito.verify;

@ExtendWith(MockitoExtension.class)
public class CreateBudgetExpenseTest {

    @Mock
    private BudgetExpenseRepository budgetExpenseRepository;

    @Mock
    private UserRepository userRepository;


    @Test
    public void saveANewBudgetExpenseWithoutASearchTag() {
        CreateBudgetExpense createBudgetExpense = new CreateBudgetExpense(budgetExpenseRepository, userRepository);

        Date date = Date.dateFor("22/10/2018");
        Money amount = Money.ONE;

        given(userRepository.currentLoggedUserName())
                .willReturn(new UserName("USER"));


        BudgetExpense budgetExpense = new BudgetExpense(null, new UserName("USER"), date, amount, "", SearchTag.DEFAULT_KEY);

        createBudgetExpense.newBudgetExpense(new NewBudgetExpenseRequest(date, amount, "", ""));

        verify(userRepository).currentLoggedUserName();
        verify(budgetExpenseRepository).save(budgetExpense);
    }
}
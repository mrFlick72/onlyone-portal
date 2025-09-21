package com.onlyoneportal.budget.expense.action;

import com.onlyoneportal.budget.Money;
import com.onlyoneportal.budget.expense.model.BudgetExpense;
import com.onlyoneportal.budget.expense.model.BudgetExpenseId;
import com.onlyoneportal.budget.expense.repository.BudgetExpenseRepository;
import com.onlyoneportal.budget.searchtag.SearchTag;
import com.onlyoneportal.budget.support.BudgetFixture;
import com.onlyoneportal.budget.time.Date;
import com.onlyoneportal.budget.user.UserName;
import org.junit.jupiter.api.Test;
import org.junit.jupiter.api.extension.ExtendWith;
import org.mockito.Mock;
import org.mockito.junit.jupiter.MockitoExtension;

import java.util.Optional;

import static org.mockito.BDDMockito.given;
import static org.mockito.Mockito.verify;

@ExtendWith(MockitoExtension.class)
public class UpdateBudgetExpenseDetailsTest {


    @Mock
    private BudgetExpenseRepository budgetExpenseRepository;

    @Test
    public void update() {
        BudgetExpenseId budgetExpenseId = BudgetFixture.randomBudgetExpenseId();

        UpdateBudgetExpenseDetails updateBudgetExpenseDetails = new UpdateBudgetExpenseDetails(budgetExpenseRepository);

        BudgetExpense budgetExpense = new BudgetExpense(budgetExpenseId, new UserName("USER"), Date.dateFor("22/02/2018"), Money.ONE, "test", SearchTag.DEFAULT_KEY);
        BudgetExpense foundBudgetExpense = new BudgetExpense(budgetExpenseId, new UserName("USER"), Date.dateFor("22/02/2018"), Money.ONE, "", SearchTag.DEFAULT_KEY);

        BudgetExpense updatedBudgetExpense = new BudgetExpense(budgetExpenseId, new UserName("USER"), Date.dateFor("22/02/2018"), Money.ONE, "test", SearchTag.DEFAULT_KEY);

        given(budgetExpenseRepository.findFor(budgetExpense.id()))
                .willReturn(Optional.of(foundBudgetExpense));

        updateBudgetExpenseDetails.updateWithoutAttachment(budgetExpense);

        verify(budgetExpenseRepository).findFor(budgetExpense.id());
        verify(budgetExpenseRepository).save(updatedBudgetExpense);
    }
}
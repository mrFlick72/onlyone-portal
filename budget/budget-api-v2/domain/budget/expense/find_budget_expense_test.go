package expense

import (
	"testing"
)

func TestBudgetExpenseTotalBySearchTagsWithConstraints(t *testing.T) {

/*
        given(userRepository.currentLoggedUserName())
                .willReturn(A_USER_NAME);

        lenient().when(searchTagRepository.findSearchTagBy("super-market")).thenReturn(new SearchTag("super-market", "super-market"));
        given(searchTagRepository.findSearchTagBy("dinner")).willReturn(new SearchTag("dinner", "dinner"));

        given(budgetExpenseRepository.findByDateRange(A_USER_NAME, Date.firstDateOfMonth(Month.FEBRUARY, Year.of(2018)),
                Date.lastDateOfMonth(Month.FEBRUARY, Year.of(2018)), "dinner", "super-market"))

                .willReturn(asList(new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), A_USER_NAME, Date.dateFor("15/02/2018"), Money.moneyFor("10"), "dinner", "dinner"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), A_USER_NAME, Date.dateFor("01/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), A_USER_NAME, Date.dateFor("05/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), A_USER_NAME, Date.dateFor("04/02/2018"), Money.moneyFor("20"), "dinner", "dinner"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), A_USER_NAME, Date.dateFor("03/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), A_USER_NAME, Date.dateFor("02/02/2018"), Money.moneyFor("12.50"), "super-market", "super-market"),
                        new BudgetExpense(BudgetFixture.emptyBudgetExpenseId(), A_USER_NAME, Date.dateFor("01/02/2018"), Money.moneyFor("15"), "dinner", "dinner")));

        SpentBudget actual = new FindSpentBudget(userRepository, budgetExpenseRepository, searchTagRepository)
                .findBy(Month.FEBRUARY, Year.of(2018), asList("dinner", "super-market"));

        Map<SearchTag, Money> expected = Map.of(new SearchTag("super-market", "super-market"), Money.moneyFor("50.00"),
                new SearchTag("dinner", "dinner"), Money.moneyFor("45.00"));

        verify(userRepository).currentLoggedUserName();
        verify(searchTagRepository).findSearchTagBy("dinner");
        verify(searchTagRepository).findSearchTagBy("super-market");
        verify(budgetExpenseRepository)
                .findByDateRange(A_USER_NAME, Date.firstDateOfMonth(Month.FEBRUARY, Year.of(2018)),
                        Date.lastDateOfMonth(Month.FEBRUARY, Year.of(2018)), "dinner", "super-market");

        Assertions.assertEquals(expected, actual.totalForSearchTags());
        Assertions.assertEquals(actual.total(), Money.moneyFor("95.00"));
		*/


}

func TestBudgetExpenseTotalBySearchTagsWithoutConstraints(t *testing.T) {
}

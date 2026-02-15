package web

import (
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)


// todo to be tested
func RepresentationModelToDomainModel(budgetExpenseRepresentation BudgetExpenseRepresentation) *expense.BudgetExpense {
	date, _ := date.DateFor(budgetExpenseRepresentation.Date)
	money, _ := money.MoneyFor(budgetExpenseRepresentation.Amount)
	return &expense.BudgetExpense{
		Id:       expense.BudgetExpenseId(budgetExpenseRepresentation.Id),
		UserName: "", // This will be set in the service layer using the current logged user
		Date:     *date,
		Amount:   money,
		Note:     budgetExpenseRepresentation.Note,
		Tag:      budgetExpenseRepresentation.TagKey,
	}
}

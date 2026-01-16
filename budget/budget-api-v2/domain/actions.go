package domain

import (
	"context"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type CreateBudgetExpense struct {
	repository BudgetExpenseRepository
}

func (action *CreateBudgetExpense) Execute(ctx *context.Context, budgetExpense *BudgetExpense) error {
	user, _ := security.GetCurrentUser(ctx)
	budgetExpense.UserName = *user.UserName
	return action.repository.Save(ctx, budgetExpense)
}

type FindBudgetExpense struct {
	repository BudgetExpenseRepository
}

type UpdateBudgetExpense struct {
	repository BudgetExpenseRepository
}

type DeleteBudgetExpense struct {
	repository BudgetExpenseRepository
}

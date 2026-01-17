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

func (action *FindBudgetExpense) Execute(ctx *context.Context, id BudgetExpenseId) (*BudgetExpense, error) {
	return nil, nil
}

type UpdateBudgetExpense struct {
	repository BudgetExpenseRepository
}

//todo add check for user ownership
func (action *UpdateBudgetExpense) Execute(ctx *context.Context, budgetExpense *BudgetExpense) error {
	existingBudgetExpense, err := action.repository.FindFor(ctx, budgetExpense.Id)
	if err != nil {
		return err
	}
	if existingBudgetExpense != nil {
		return action.repository.Save(ctx, budgetExpense)
	}
	return nil
}

type DeleteBudgetExpense struct {
	repository BudgetExpenseRepository
}

func (action *DeleteBudgetExpense) Execute(ctx *context.Context, id BudgetExpenseId) error {
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingBudgetExpense, err := action.repository.FindFor(ctx, id)
	if err != nil {
		return err
	}

	if existingBudgetExpense != nil && existingBudgetExpense.UserName == *userName.UserName {
		return action.repository.Delete(ctx, id)
	}
	return nil
}

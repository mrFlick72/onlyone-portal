package expense

import (
	"context"
	"errors"

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

func (action *FindBudgetExpense) ExecuteAll(ctx *context.Context) (*SpentBudget, error) {
	return nil, nil
}

func (action *FindBudgetExpense) Execute(ctx *context.Context, id BudgetExpenseId) (*BudgetExpense, error) {
	return nil, nil
}

type UpdateBudgetExpense struct {
	repository BudgetExpenseRepository
}

func (action *UpdateBudgetExpense) Execute(ctx *context.Context, budgetExpense *BudgetExpense) error {
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingBudgetExpense, err := action.repository.FindFor(ctx, budgetExpense.Id)
	if err != nil {
		return err
	}
	if existingBudgetExpense != nil && existingBudgetExpense.UserName == *userName.UserName {
		return action.repository.Save(ctx, budgetExpense)
	}
	return errors.New("budget expense not found or user not authorized to delete it")
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
	return errors.New("budget expense not found or user not authorized to delete it")
}

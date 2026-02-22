package expense

import (
	"context"

	"github.com/mrflick72/budget/budget-api/adapter/tags/rest"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
)

func NewBudgetExpenseActionsFacade() BudgetExpenseActions {
	budgetExpenseRepository := NewBudgetExpenseRepository()
	searchTagRepository := rest.NewSearchTagRepository()
	createBudgetExpense := &CreateBudgetExpense{
		repository: budgetExpenseRepository,
	}
	updateBudgetExpense := &UpdateBudgetExpense{
		repository: budgetExpenseRepository,
	}
	findSpentBudget := &FindSpentBudget{
		budgetExpenseRepository: budgetExpenseRepository,
		searchTagRepository:     searchTagRepository,
	}
	deleteBudgetExpense := &DeleteBudgetExpense{
		repository: budgetExpenseRepository,
	}

	return &BudgetExpenseActionsFacade{
		CreateBudgetExpenseAction: createBudgetExpense,
		UpdateBudgetExpenseAction: updateBudgetExpense,
		FindSpentBudgetAction:     findSpentBudget,
		DeleteBudgetExpenseAction: deleteBudgetExpense,
	}
}

func NewBudgetExpenseRepository() BudgetExpenseRepository {
	panic("unimplemented")
}

type BudgetExpenseActions interface {
	CreateBudgetExpense(ctx *context.Context, budgetExpense *BudgetExpense) error
	UpdateBudgetExpense(ctx *context.Context, budgetExpense *BudgetExpense) error
	FindSpentBudget(ctx *context.Context, month date.Month, year date.Year, searchTagKeys []tags.SearchTagKey) (*SpentBudget, error)
	DeleteBudgetExpense(ctx *context.Context, id BudgetExpenseId) error
}

type BudgetExpenseActionsFacade struct {
	CreateBudgetExpenseAction *CreateBudgetExpense
	UpdateBudgetExpenseAction *UpdateBudgetExpense
	FindSpentBudgetAction     *FindSpentBudget
	DeleteBudgetExpenseAction *DeleteBudgetExpense
}

func (facade *BudgetExpenseActionsFacade) CreateBudgetExpense(ctx *context.Context, budgetExpense *BudgetExpense) error {
	return facade.CreateBudgetExpenseAction.Execute(ctx, budgetExpense)
}

func (facade *BudgetExpenseActionsFacade) UpdateBudgetExpense(ctx *context.Context, budgetExpense *BudgetExpense) error {
	return facade.UpdateBudgetExpenseAction.Execute(ctx, budgetExpense)
}

func (facade *BudgetExpenseActionsFacade) FindSpentBudget(ctx *context.Context, month date.Month, year date.Year, searchTagKeys []tags.SearchTagKey) (*SpentBudget, error) {
	return facade.FindSpentBudgetAction.Execute(ctx, month, year, searchTagKeys)
}

func (facade *BudgetExpenseActionsFacade) DeleteBudgetExpense(ctx *context.Context, id BudgetExpenseId) error {
	return facade.DeleteBudgetExpenseAction.Execute(ctx, id)
}

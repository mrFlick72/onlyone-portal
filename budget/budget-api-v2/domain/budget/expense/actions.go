package expense

import (
	"context"
	"errors"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type BudgetExpenseActions interface {
	CreateBudgetExpense(ctx *context.Context, budgetExpense *BudgetExpense) error
	UpdateBudgetExpense(ctx *context.Context, budgetExpense *BudgetExpense) error
	FindSpentBudget(ctx *context.Context, month date.Month, year date.Year, searchTagKeys []tags.SearchTagKey) (*SpentBudget, error)
	DeleteBudgetExpense(ctx *context.Context, id BudgetExpenseId) error
}

type BudgetExpenseActionsFacade struct {
	CreateBudgetExpense *CreateBudgetExpense
	UpdateBudgetExpense *UpdateBudgetExpense
	FindSpentBudget     *FindSpentBudget
	DeleteBudgetExpense *DeleteBudgetExpense
}

type CreateBudgetExpense struct {
	repository BudgetExpenseRepository
}

func (action *CreateBudgetExpense) Execute(ctx *context.Context, budgetExpense *BudgetExpense) error {
	user, _ := security.GetCurrentUser(ctx)
	budgetExpense.UserName = *user.UserName
	return action.repository.Save(ctx, budgetExpense)
}

type FindSpentBudget struct {
	budgetExpenseRepository BudgetExpenseRepository
	searchTagRepository     tags.SearchTagRepository
}

func (action *FindSpentBudget) Execute(ctx *context.Context, month date.Month, year date.Year, searchTagKeys []tags.SearchTagKey) (*SpentBudget, error) {
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	firstDate, err := date.FirstDateOfMonth(month, year)
	if err != nil {
		return nil, err
	}
	lastDate, err := date.LastDateOfMonth(month, year)
	if err != nil {
		return nil, err
	}

	budgetByDateRange, err := action.budgetExpenseRepository.FindByDateRange(ctx, *userName.UserName, *firstDate, *lastDate, searchTagKeys)
	if err != nil {
		return nil, err
	}

	searchTags, err := action.getAllSearchTagFor(ctx, budgetByDateRange)
	return NewSpentBudget(budgetByDateRange, searchTags), err
}

func (action *FindSpentBudget) getAllSearchTagFor(ctx *context.Context, budgetExpenses *[]BudgetExpense) (*[]tags.SearchTag, error) {
	var searchTags []tags.SearchTag
	seen := make(map[string]bool)
	for _, expense := range *budgetExpenses {
		if !seen[expense.Tag] {
			seen[expense.Tag] = true
			searchTag, err := action.searchTagRepository.GetTagBy(ctx, expense.Tag)
			if err != nil {
				return nil, err
			}
			searchTags = append(searchTags, *searchTag)
		}
	}
	return &searchTags, nil
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

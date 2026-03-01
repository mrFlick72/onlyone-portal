package expense

import (
	"context"
	"errors"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)


type CreateBudgetExpense struct {
	Repository BudgetExpenseRepository
}

func (action *CreateBudgetExpense) Execute(ctx context.Context, budgetExpense *BudgetExpense) error {
	user, _ := security.GetCurrentUser(ctx)
	budgetExpense.UserName = *user.UserName
	return action.Repository.Save(ctx, budgetExpense)
}

type FindSpentBudget struct {
	BudgetExpenseRepository BudgetExpenseRepository
	SearchTagRepository     tags.SearchTagRepository
}

func (action *FindSpentBudget) Execute(ctx context.Context, month date.Month, year date.Year, searchTagKeys []tags.SearchTagKey) (*SpentBudget, error) {
	firstDate, err := date.FirstDateOfMonth(month, year)
	if err != nil {
		return nil, err
	}
	lastDate, err := date.LastDateOfMonth(month, year)
	if err != nil {
		return nil, err
	}

	budgetByDateRange, err := action.BudgetExpenseRepository.FindByDateRange(ctx, firstDate, lastDate, searchTagKeys)
	if err != nil {
		return nil, err
	}

	searchTags, err := action.getAllSearchTagFor(ctx, budgetByDateRange)
	return NewSpentBudget(budgetByDateRange, searchTags), err
}

func (action *FindSpentBudget) getAllSearchTagFor(ctx context.Context, budgetExpenses []BudgetExpense) ([]tags.SearchTag, error) {
	var searchTags []tags.SearchTag
	seen := make(map[string]bool)
	for _, expense := range budgetExpenses {
		if !seen[expense.Tag.Key] {
			seen[expense.Tag.Key] = true
			searchTag, err := action.SearchTagRepository.GetTagBy(ctx, expense.Tag.Key)
			if err != nil {
				return nil, err
			}
			searchTags = append(searchTags, *searchTag)
		}
	}
	return searchTags, nil
}

type UpdateBudgetExpense struct {
	Repository BudgetExpenseRepository
}

func (action *UpdateBudgetExpense) Execute(ctx context.Context, budgetExpense *BudgetExpense) error {
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingBudgetExpense, err := action.Repository.FindFor(ctx, budgetExpense.Id)
	if err != nil {
		return err
	}
	if existingBudgetExpense != nil && existingBudgetExpense.UserName == *userName.UserName {
		return action.Repository.Save(ctx, budgetExpense)
	}
	return errors.New("budget expense not found or user not authorized to update it")
}

type DeleteBudgetExpense struct {
	Repository BudgetExpenseRepository
}

func (action *DeleteBudgetExpense) Execute(ctx context.Context, id BudgetExpenseId) error {
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingBudgetExpense, err := action.Repository.FindFor(ctx, id)
	if err != nil {
		return err
	}

	if existingBudgetExpense != nil && existingBudgetExpense.UserName == *userName.UserName {
		return action.Repository.Delete(ctx, id)
	}
	return errors.New("budget expense not found or user not authorized to delete it")
}

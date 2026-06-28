package expense

import (
	"context"
	"errors"
	"fmt"

	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type CreateBudgetExpense struct {
	Repository     BudgetExpenseRepository
	EventPublisher BudgetExpenseEventPublisher
}

func (action *CreateBudgetExpense) Execute(ctx context.Context, budgetExpense *BudgetExpense) error {
	applyDefaultTagIfMissing(budgetExpense)
	user, _ := security.GetCurrentUser(ctx)
	budgetExpense.UserName = *user.UserName
	err := action.Repository.Save(ctx, budgetExpense)
	if err == nil {
		action.EventPublisher.CreateBudgetExpense(ctx, *budgetExpense)
	}
	return err
}

// applyDefaultTagIfMissing defaults an expense with no tags to the UNKNOWN sentinel tag,
// which tag-api always includes in its catalog without requiring per-user onboarding
// (see tagging/tag-api/docs/adr/0001-unknown-sentinel-tag-not-persisted.md). The
// sentinel is defined once in the tags package and shared with revenue.
func applyDefaultTagIfMissing(budgetExpense *BudgetExpense) {
	if len(budgetExpense.Tags) == 0 {
		budgetExpense.Tags = []tags.SearchTag{tags.UnknownSentinel()}
	}
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
		for _, tag := range expense.Tags {
			if !seen[tag.Key] {
				seen[tag.Key] = true
				searchTag, err := action.SearchTagRepository.GetTagBy(ctx, tag.Key)
				if err != nil {
					return nil, err
				}
				searchTags = append(searchTags, *searchTag)
			}
		}

	}
	return searchTags, nil
}

type UpdateBudgetExpense struct {
	Repository     BudgetExpenseRepository
	EventPublisher BudgetExpenseEventPublisher
	EventBus       EventBus
}

func (action *UpdateBudgetExpense) Listen() {
	for {
		event := <-action.EventBus

		fmt.Println(event)
		fmt.Println("event listener event consumed")
	}
}

func (action *UpdateBudgetExpense) Execute(ctx context.Context, budgetExpense *BudgetExpense) error {
	applyDefaultTagIfMissing(budgetExpense)
	userName, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	existingBudgetExpense, err := action.Repository.FindFor(ctx, budgetExpense.Id)
	if err != nil {
		return err
	}
	if existingBudgetExpense != nil && existingBudgetExpense.UserName == *userName.UserName {
		err := action.Repository.Save(ctx, budgetExpense)
		if err == nil {
			action.EventPublisher.UpdateBudgetExpense(ctx, *budgetExpense)
		}
		return err
	}
	return errors.New("budget expense not found or user not authorized to update it")
}

type DeleteBudgetExpense struct {
	Repository     BudgetExpenseRepository
	EventPublisher BudgetExpenseEventPublisher
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
		err := action.Repository.Delete(ctx, id)
		if err == nil {
			action.EventPublisher.DeleteBudgetExpense(ctx, *existingBudgetExpense)
		}
		return err
	}
	return errors.New("budget expense not found or user not authorized to delete it")
}

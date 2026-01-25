package dynamodb

import (
	"context"

	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type DynamoDbBudgetExpenseRepository struct {
}

func (repository *DynamoDbBudgetExpenseRepository) FindByDateRange(ctx *context.Context, userName security.UserName, start date.Date, end date.Date, searchTags []tags.SearchTagKey) (*[]expense.BudgetExpense, error) {
	// Implementation to interact with DynamoDB and retrieve budget expenses by date range and search tags
	return nil, nil
}

func (repository *DynamoDbBudgetExpenseRepository) Save(ctx *context.Context, budgetExpense *expense.BudgetExpense) error {
	// Implementation to save a budget expense to DynamoDB
	return nil
}

func (repository *DynamoDbBudgetExpenseRepository) Delete(ctx *context.Context, idBudgetExpense expense.BudgetExpenseId) error {
	// Implementation to delete a budget expense from DynamoDB
	return nil
}

func (repository *DynamoDbBudgetExpenseRepository) FindFor(ctx *context.Context, budgetExpenseId expense.BudgetExpenseId) (*expense.BudgetExpense, error) {
	// Implementation to find a budget expense by its ID in DynamoDB
	return nil, nil
}

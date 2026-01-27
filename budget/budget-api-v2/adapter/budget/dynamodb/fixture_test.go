package dynamodb

import (
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/stretchr/testify/mock"
)

type DynamoDbBudgetExpenseIdProviderMock struct {
	mock.Mock
}

func (mock *DynamoDbBudgetExpenseIdProviderMock) GenerateIdFor(budgetExpense *expense.BudgetExpense) string {
	args := mock.Called(budgetExpense)
	return args.String(0)
}

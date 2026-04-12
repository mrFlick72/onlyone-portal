package dynamodb

import (
	"encoding/base64"
	"fmt"

	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type DynamoDbBudgetExpenseIdProvider struct {
	SaltGenerator func() string 
	
}

func (provider *DynamoDbBudgetExpenseIdProvider) GenerateIdFor(budgetExpense *expense.BudgetExpense) expense.BudgetExpenseId {
	return fmt.Sprintf("%s-%s", 
	provider.partitionKeyFrom(budgetExpense.Date.GetYear(), budgetExpense.Date.GetMonth(), budgetExpense.UserName), 
	provider.rangeKeyFrom(budgetExpense.Date.GetDay()))
}

func (provider *DynamoDbBudgetExpenseIdProvider) partitionKeyFrom(year, month int, userName security.UserName) string {
	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d_%d_%s", year, month, userName)))
}

func (provider *DynamoDbBudgetExpenseIdProvider) rangeKeyFrom(day int) string {
	rangeKey := fmt.Sprintf("%d_%s", day, provider.SaltGenerator())
	return base64.StdEncoding.EncodeToString([]byte(rangeKey))
}

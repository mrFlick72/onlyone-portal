package dynamodb

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type DynamoDbBudgetExpenseIdProvider struct {
	saltGenerator func() string
}

func (provider *DynamoDbBudgetExpenseIdProvider) GenerateIdFor(budgetExpense *expense.BudgetExpense) expense.BudgetExpenseId {
	return fmt.Sprintf("%s-%s", provider.partitionKeyFrom(budgetExpense.Date, budgetExpense.UserName), provider.rangeKeyFrom(budgetExpense))
}

func (provider *DynamoDbBudgetExpenseIdProvider) partitionKeyFrom(date date.Date, userName security.UserName) string {
	isoDate := date.GetIsoFormattedDate()
	budgetExpenseYear, _ := strconv.Atoi(isoDate[:4])
	budgetExpenseMonth, _ := strconv.Atoi(isoDate[5:7])

	return base64.StdEncoding.EncodeToString([]byte(fmt.Sprintf("%d_%d_%s", budgetExpenseYear, budgetExpenseMonth, userName)))
}

func (provider *DynamoDbBudgetExpenseIdProvider) rangeKeyFrom(budgetExpense *expense.BudgetExpense) string {
	isoDate := budgetExpense.Date.GetIsoFormattedDate()
	budgetExpenseDay, _ := strconv.Atoi(isoDate[8:10])
	rangeKey := fmt.Sprintf("%d_%s", budgetExpenseDay, provider.saltGenerator())
	return base64.StdEncoding.EncodeToString([]byte(rangeKey))
}

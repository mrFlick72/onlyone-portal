package dynamodb

import (
	"testing"

	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func TestANewBudgetIdCreation(t *testing.T) {
	provider := &DynamoDbBudgetExpenseIdProvider{
		SaltGenerator: func() string { return "A_SALT" },
	}
	budgetExpense := expense.BudgetExpense{
		UserName: "USER",
		Date:     testutils.SafeDateFor("01/01/2018"),
		Amount:   testutils.SafeMoneyFor("1.00"),
		Note:     "",
		Tag:      tags.SearchTag{Key: "TAG", Value: "TAG"},
	}

	id := provider.GenerateIdFor(&budgetExpense)
	assert.Equal(t, "MjAxOF8xX1VTRVI=-MV9BX1NBTFQ=", id)
}

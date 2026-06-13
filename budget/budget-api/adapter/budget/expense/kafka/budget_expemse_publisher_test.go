package kafka

import (
	"encoding/json"
	"testing"

	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/stretchr/testify/assert"
)

func aBudgetExpense() expense.BudgetExpense {
	return expense.BudgetExpense{
		Id:       "A_BUDGET_ID",
		UserName: "testuser",
		Date:     testutils.SafeDateFor("01/02/2024"),
		Amount:   testutils.SafeMoneyFor("10.50"),
		Note:     "lunch with the team",
		Tags:     []tags.SearchTag{{Key: "lunch", Value: "food"}},
	}
}

func TestNewBudgetExpenseEventMapsAllFields(t *testing.T) {
	event := newBudgetExpenseEvent(CreateAction, aBudgetExpense())

	assert.Equal(t, CreateAction, event.Action)
	assert.Equal(t, BudgetExpenseBody{
		Id:       "A_BUDGET_ID",
		UserName: "testuser",
		Date:     "01/02/2024",
		Amount:   "10.50",
		Note:     "lunch with the team",
		Tags:     []tags.SearchTag{{Key: "lunch", Value: "food"}},
	}, event.Payload)
}

func TestEncodeBudgetExpenseEventProducesNonEmptyPayload(t *testing.T) {
	encoded, err := encodeBudgetExpenseEvent(CreateAction, aBudgetExpense())
	assert.NoError(t, err)

	var decoded BudgetExpenseEvent
	assert.NoError(t, json.Unmarshal(encoded, &decoded))

	// Guards against the regression where unexported fields (or the raw domain
	// Money/Date value objects) serialize to an empty {} and drop the amount,
	// the date and the action.
	assert.Equal(t, CreateAction, decoded.Action)
	assert.Equal(t, "10.50", decoded.Payload.Amount)
	assert.Equal(t, "01/02/2024", decoded.Payload.Date)
	assert.Equal(t, "testuser", decoded.Payload.UserName)
	assert.Equal(t, "lunch with the team", decoded.Payload.Note)
	assert.Equal(t, []tags.SearchTag{{Key: "lunch", Value: "food"}}, decoded.Payload.Tags)
}

func TestEncodeBudgetExpenseEventJSONShape(t *testing.T) {
	encoded, err := encodeBudgetExpenseEvent(UpdateAction, aBudgetExpense())
	assert.NoError(t, err)

	assert.JSONEq(t, `{
		"action": "UPDATE",
		"payload": {
			"id": "A_BUDGET_ID",
			"userName": "testuser",
			"date": "01/02/2024",
			"amount": "10.50",
			"note": "lunch with the team",
			"tags": [{"key": "lunch", "value": "food"}]
		}
	}`, string(encoded))
}

func TestNewBudgetExpenseEventActionMapping(t *testing.T) {
	for _, action := range []string{CreateAction, UpdateAction, DeleteAction} {
		t.Run(action, func(t *testing.T) {
			event := newBudgetExpenseEvent(action, aBudgetExpense())
			assert.Equal(t, action, event.Action)
		})
	}
}

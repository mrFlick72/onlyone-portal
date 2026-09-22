//go:build test

package dynamodb

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
)

func TestGenerateIdForReturnsExistingWhenPresent(t *testing.T) {
	provider := &DynamoDbScheduledExpenseIdProvider{
		UuidGenerator: func() string { return "GENERATED" },
	}

	se := &scheduledexpense.ScheduledExpense{Id: "EXISTING"}

	assert.Equal(t, "EXISTING", provider.GenerateIdFor(se))
}

func TestGenerateIdForGeneratesWhenMissing(t *testing.T) {
	provider := &DynamoDbScheduledExpenseIdProvider{
		UuidGenerator: func() string { return "GENERATED" },
	}

	se := &scheduledexpense.ScheduledExpense{Id: ""}

	assert.Equal(t, "GENERATED", provider.GenerateIdFor(se))
}

//go:build test

package dynamodb

import (
	"testing"

	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
)

func TestPartitionKeyForExpenseUppercasesType(t *testing.T) {
	provider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "" },
	}

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:   "budget-123",
			BudgetType: "expense",
		},
	}

	assert.Equal(t, "budget-123_EXPENSE", provider.PartitionKeyFor(att))
}

func TestPartitionKeyForRevenueUppercasesType(t *testing.T) {
	provider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "" },
	}

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:   "rev-9",
			BudgetType: "revenue",
		},
	}

	assert.Equal(t, "rev-9_REVENUE", provider.PartitionKeyFor(att))
}

func TestGenerateAttachmentIdReturnsExistingWhenPresent(t *testing.T) {
	provider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "GENERATED" },
	}

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{AttachmentId: "EXISTING"},
	}

	assert.Equal(t, "EXISTING", provider.GenerateAttachmentIdFor(att))
}

func TestGenerateAttachmentIdGeneratesWhenMissing(t *testing.T) {
	provider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "GENERATED" },
	}

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{AttachmentId: ""},
	}

	assert.Equal(t, "GENERATED", provider.GenerateAttachmentIdFor(att))
}

func TestRangeKeyForReturnsAttachmentId(t *testing.T) {
	provider := &DynamoDbAttachmentIdProvider{}
	assert.Equal(t, "ATT-1", provider.RangeKeyFor("ATT-1"))
}

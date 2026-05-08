//go:build test

package dynamodb

import (
	"context"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/go-playground/assert/v2"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

func TestMain(m *testing.M) {
	setupTestDynamoDBTable()
	code := m.Run()
	teardownTestDynamoDBTable()
	os.Exit(code)
}

func TestSaveAttachmentMetadataWritesAllFields(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-001",
			BudgetId:     "budget-123",
			BudgetType:   "expense",
			Date:         testutils.SafeDateFor("15/03/2024"),
			Owner:        "testuser",
			FineName:     "receipt.pdf",
			ContentType:  "application/pdf",
			Metadata: map[string]string{
				"source": "mobile",
			},
		},
	}
	pk := "budget-123_EXPENSE"
	rk := "att-001"
	fileLocation := "test-bucket/2024/03/15/budget-123_EXPENSE/att-001"

	err := repo.Save(context.Background(), att, pk, rk, fileLocation)
	assert.Equal(t, nil, err)

	item, err := getItem(context.Background(), pk, rk)
	assert.Equal(t, nil, err)

	assert.Equal(t, "att-001", item["attachment_id"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "budget-123", item["budget_id"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "expense", item["budget_type"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "2024-03-15", item["date"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "testuser", item["user_name"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "receipt.pdf", item["file_name"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, "application/pdf", item["content_type"].(*types.AttributeValueMemberS).Value)
	assert.Equal(t, fileLocation, item["file_location"].(*types.AttributeValueMemberS).Value)

	metadataAttr := item["metadata"].(*types.AttributeValueMemberM).Value
	assert.Equal(t, "mobile", metadataAttr["source"].(*types.AttributeValueMemberS).Value)
}

func TestSaveAttachmentMetadataOmitsEmptyMetadataMap(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-002",
			BudgetId:     "budget-456",
			BudgetType:   "revenue",
			Date:         testutils.SafeDateFor("01/06/2024"),
			Owner:        "testuser",
			FineName:     "invoice.png",
			ContentType:  "image/png",
		},
	}
	pk := "budget-456_REVENUE"
	rk := "att-002"

	err := repo.Save(context.Background(), att, pk, rk, "test-bucket/2024/06/01/budget-456_REVENUE/att-002")
	assert.Equal(t, nil, err)

	item, err := getItem(context.Background(), pk, rk)
	assert.Equal(t, nil, err)

	_, hasMetadata := item["metadata"]
	assert.Equal(t, false, hasMetadata)
}

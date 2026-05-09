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
	fileLocation := "test-bucket/2024/03/15/budget-123_EXPENSE/att-001"

	err := repo.Save(context.Background(), att, pk, fileLocation)
	assert.Equal(t, nil, err)

	item, err := getItem(context.Background(), pk, att.AttachmentId)
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

	err := repo.Save(context.Background(), att, pk, "test-bucket/2024/06/01/budget-456_REVENUE/att-002")
	assert.Equal(t, nil, err)

	item, err := getItem(context.Background(), pk, att.AttachmentId)
	assert.Equal(t, nil, err)

	_, hasMetadata := item["metadata"]
	assert.Equal(t, false, hasMetadata)
}

func TestFindAllAttachmentReturnsItemsForOwnerAndPartition(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	pk := "budget-find_EXPENSE"
	att1 := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-find-1",
			BudgetId:     "budget-find",
			BudgetType:   "expense",
			Date:         testutils.SafeDateFor("10/01/2024"),
			Owner:        "owner-alice",
			FineName:     "alice-1.pdf",
			ContentType:  "application/pdf",
			Metadata: map[string]string{
				attachment.MetadataKeyBucket:    "alice-bucket",
				attachment.MetadataKeyObjectKey: "2024/01/10/budget-find_EXPENSE/att-find-1",
				"source":                        "web",
			},
		},
	}
	att2 := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-find-2",
			BudgetId:     "budget-find",
			BudgetType:   "expense",
			Date:         testutils.SafeDateFor("11/01/2024"),
			Owner:        "owner-bob",
			FineName:     "bob.pdf",
			ContentType:  "application/pdf",
		},
	}

	assert.Equal(t, nil, repo.Save(context.Background(), att1, pk, "alice-bucket/2024/01/10/budget-find_EXPENSE/att-find-1"))
	assert.Equal(t, nil, repo.Save(context.Background(), att2, pk, "alice-bucket/2024/01/11/budget-find_EXPENSE/att-find-2"))

	got, err := repo.FindAllAttachment(context.Background(), "owner-alice", pk)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, "att-find-1", got[0].AttachmentId)
	assert.Equal(t, "owner-alice", got[0].Owner)
	assert.Equal(t, "alice-1.pdf", got[0].FineName)
	assert.Equal(t, "alice-bucket", got[0].Metadata[attachment.MetadataKeyBucket])
	assert.Equal(t, "2024/01/10/budget-find_EXPENSE/att-find-1", got[0].Metadata[attachment.MetadataKeyObjectKey])
	assert.Equal(t, "web", got[0].Metadata["source"])
}

func TestFindAllAttachmentReturnsEmptySliceWhenNoMatches(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	got, err := repo.FindAllAttachment(context.Background(), "no-owner", "no_PARTITION")
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(got))
}

func TestGetAttachmentReturnsAttachmentWithMetadataMap(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	pk := "budget-get_EXPENSE"
	objectKey := "2024/02/02/budget-get_EXPENSE/att-get-1"
	fileLocation := "get-bucket/" + objectKey
	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-get-1",
			BudgetId:     "budget-get",
			BudgetType:   "expense",
			Date:         testutils.SafeDateFor("02/02/2024"),
			Owner:        "owner-get",
			FineName:     "doc.pdf",
			ContentType:  "application/pdf",
			Metadata: map[string]string{
				attachment.MetadataKeyBucket:    "get-bucket",
				attachment.MetadataKeyObjectKey: objectKey,
			},
		},
	}
	assert.Equal(t, nil, repo.Save(context.Background(), att, pk, fileLocation))

	loaded, err := repo.GetAttachment(context.Background(), "owner-get", "att-get-1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "att-get-1", loaded.AttachmentId)
	assert.Equal(t, "owner-get", loaded.Owner)
	assert.Equal(t, "doc.pdf", loaded.FineName)
	assert.Equal(t, "application/pdf", loaded.ContentType)
	assert.Equal(t, "get-bucket", loaded.Metadata[attachment.MetadataKeyBucket])
	assert.Equal(t, objectKey, loaded.Metadata[attachment.MetadataKeyObjectKey])
	assert.Equal(t, fileLocation, loaded.Metadata["file_location"])
}

func TestGetAttachmentReturnsErrorWhenNotFound(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	_, err := repo.GetAttachment(context.Background(), "any-owner", "missing-id")
	assert.NotEqual(t, nil, err)
}

func TestDeleteReturnsObjectKeyAndRemovesItem(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	pk := "budget-del_EXPENSE"
	objectKey := "2024/03/03/budget-del_EXPENSE/att-del-1"
	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-del-1",
			BudgetId:     "budget-del",
			BudgetType:   "expense",
			Date:         testutils.SafeDateFor("03/03/2024"),
			Owner:        "owner-del",
			FineName:     "del.pdf",
			ContentType:  "application/pdf",
			Metadata: map[string]string{
				attachment.MetadataKeyBucket:    "del-bucket",
				attachment.MetadataKeyObjectKey: objectKey,
			},
		},
	}
	assert.Equal(t, nil, repo.Save(context.Background(), att, pk, "del-bucket/"+objectKey))

	got, err := repo.Delete(context.Background(), "owner-del", "att-del-1")
	assert.Equal(t, nil, err)
	assert.Equal(t, objectKey, got)

	item, err := getItem(context.Background(), pk, "att-del-1")
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(item))
}

func TestDeleteReturnsErrorWhenNotFound(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	_, err := repo.Delete(context.Background(), "any-owner", "missing-id")
	assert.NotEqual(t, nil, err)
}

func TestDeleteRefusesToDeleteWhenOwnerDoesNotMatch(t *testing.T) {
	idProvider := &DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return "ignored-here" },
	}
	repo := NewDynamoDbAttachmentMetadataRepository(TableName, client, idProvider)

	pk := "budget-acl_EXPENSE"
	objectKey := "2024/04/04/budget-acl_EXPENSE/att-acl-1"
	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "att-acl-1",
			BudgetId:     "budget-acl",
			BudgetType:   "expense",
			Date:         testutils.SafeDateFor("04/04/2024"),
			Owner:        "owner-real",
			FineName:     "acl.pdf",
			ContentType:  "application/pdf",
			Metadata: map[string]string{
				attachment.MetadataKeyObjectKey: objectKey,
			},
		},
	}
	assert.Equal(t, nil, repo.Save(context.Background(), att, pk, "acl-bucket/"+objectKey))

	_, err := repo.Delete(context.Background(), "intruder", "att-acl-1")
	assert.NotEqual(t, nil, err)

	item, err := getItem(context.Background(), pk, "att-acl-1")
	assert.Equal(t, nil, err)
	assert.Equal(t, "att-acl-1", item["attachment_id"].(*types.AttributeValueMemberS).Value)
}

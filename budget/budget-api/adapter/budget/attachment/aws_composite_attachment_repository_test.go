//go:build test

package attachment

import (
	"context"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/go-playground/assert/v2"
	attachmentdynamo "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/dynamodb"
	attachments3 "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/s3"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
	"github.com/mrflick72/budget/budget-api/internal/testutils/awsfixture"
)

const (
	testTableName  = "BUDGET_ATTACHMENT_METADATA_COMPOSITE_TABLE"
	testBucketName = "budget-attachment-composite-staging"
)

var dynamoClient = awsfixture.NewLocalDynamoDBClient()
var s3Client = awsfixture.NewLocalS3Client()

func TestMain(m *testing.M) {
	ctx := context.TODO()
	_ = awsfixture.TeardownTable(ctx, dynamoClient, testTableName)
	_ = awsfixture.TeardownBucket(ctx, s3Client, testBucketName)
	_ = awsfixture.SetupAttachmentTable(ctx, dynamoClient, testTableName)
	_ = awsfixture.SetupBucket(ctx, s3Client, testBucketName)
	code := m.Run()
	_ = awsfixture.TeardownBucket(ctx, s3Client, testBucketName)
	_ = awsfixture.TeardownTable(ctx, dynamoClient, testTableName)
	os.Exit(code)
}

func newCompositeRepository(uuid string) *AwsCompositeAttachmentRepository {
	idProvider := &attachmentdynamo.DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return uuid },
	}
	metadataRepo := attachmentdynamo.NewDynamoDbAttachmentMetadataRepository(testTableName, dynamoClient, idProvider)
	contentRepo := attachments3.NewS3AttachmentContentRepository(testBucketName, s3Client)
	return NewAwsCompositeAttachmentRepository(idProvider, metadataRepo, contentRepo).(*AwsCompositeAttachmentRepository)
}

func TestSaveAttachmentWritesContentToS3AndMetadataToDynamo(t *testing.T) {
	repo := newCompositeRepository("generated-uuid-1")

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:    "budget-abc",
			BudgetType:  "expense",
			Date:        testutils.SafeDateFor("15/03/2024"),
			Owner:       "testuser",
			FineName:    "receipt.pdf",
			ContentType: "application/pdf",
			Metadata:    map[string]string{"source": "web"},
		},
		Content: []byte("file-bytes"),
	}

	err := repo.SaveAttachment(context.Background(), att)
	assert.Equal(t, nil, err)
	assert.Equal(t, "generated-uuid-1", att.AttachmentId)

	expectedKey := "2024/03/15/budget-abc_EXPENSE/generated-uuid-1"

	objOut, err := s3Client.GetObject(context.Background(), &aws_s3.GetObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(expectedKey),
	})
	assert.Equal(t, nil, err)
	defer objOut.Body.Close()
	body, _ := io.ReadAll(objOut.Body)
	assert.Equal(t, []byte("file-bytes"), body)
	assert.Equal(t, "application/pdf", aws.ToString(objOut.ContentType))

	itemOut, err := dynamoClient.GetItem(context.Background(), &aws_dynamodb.GetItemInput{
		TableName: aws.String(testTableName),
		Key: map[string]dynamotypes.AttributeValue{
			"pk":            &dynamotypes.AttributeValueMemberS{Value: "budget-abc_EXPENSE"},
			"attachment_id": &dynamotypes.AttributeValueMemberS{Value: "generated-uuid-1"},
		},
	})
	assert.Equal(t, nil, err)
	expectedFileLocation := testBucketName + "/" + expectedKey
	assert.Equal(t, expectedFileLocation, itemOut.Item["file_location"].(*dynamotypes.AttributeValueMemberS).Value)
	assert.Equal(t, "generated-uuid-1", itemOut.Item["attachment_id"].(*dynamotypes.AttributeValueMemberS).Value)
	assert.Equal(t, "testuser", itemOut.Item["user_name"].(*dynamotypes.AttributeValueMemberS).Value)

	metadataAttr := itemOut.Item["metadata"].(*dynamotypes.AttributeValueMemberM).Value
	assert.Equal(t, testBucketName, metadataAttr[attachment.MetadataKeyBucket].(*dynamotypes.AttributeValueMemberS).Value)
	assert.Equal(t, expectedKey, metadataAttr[attachment.MetadataKeyObjectKey].(*dynamotypes.AttributeValueMemberS).Value)
	assert.Equal(t, "web", metadataAttr["source"].(*dynamotypes.AttributeValueMemberS).Value)
}

func TestSaveAttachmentReusesExistingAttachmentIdAndOverwrites(t *testing.T) {
	repo := newCompositeRepository("should-not-be-used")

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "preexisting-id",
			BudgetId:     "budget-xyz",
			BudgetType:   "revenue",
			Date:         testutils.SafeDateFor("01/06/2024"),
			Owner:        "testuser",
			FineName:     "old.pdf",
			ContentType:  "application/pdf",
		},
		Content: []byte("v1"),
	}

	err := repo.SaveAttachment(context.Background(), att)
	assert.Equal(t, nil, err)
	assert.Equal(t, "preexisting-id", att.AttachmentId)

	att.Content = []byte("v2")
	att.FineName = "new.pdf"

	err = repo.SaveAttachment(context.Background(), att)
	assert.Equal(t, nil, err)
	assert.Equal(t, "preexisting-id", att.AttachmentId)

	expectedKey := "2024/06/01/budget-xyz_REVENUE/preexisting-id"
	objOut, err := s3Client.GetObject(context.Background(), &aws_s3.GetObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(expectedKey),
	})
	assert.Equal(t, nil, err)
	defer objOut.Body.Close()
	body, _ := io.ReadAll(objOut.Body)
	assert.Equal(t, []byte("v2"), body)

	itemOut, err := dynamoClient.GetItem(context.Background(), &aws_dynamodb.GetItemInput{
		TableName: aws.String(testTableName),
		Key: map[string]dynamotypes.AttributeValue{
			"pk":            &dynamotypes.AttributeValueMemberS{Value: "budget-xyz_REVENUE"},
			"attachment_id": &dynamotypes.AttributeValueMemberS{Value: "preexisting-id"},
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, "new.pdf", itemOut.Item["file_name"].(*dynamotypes.AttributeValueMemberS).Value)
}

func TestGenAttachmentReturnsContentRoundTrippedFromS3(t *testing.T) {
	repo := newCompositeRepository("gen-uuid")

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:    "budget-gen",
			BudgetType:  "expense",
			Date:        testutils.SafeDateFor("12/04/2024"),
			Owner:       "gen-user",
			FineName:    "gen.pdf",
			ContentType: "application/pdf",
		},
		Content: []byte("gen-bytes"),
	}
	ctx := testutils.NewStubbedContextWith("gen-user")
	assert.Equal(t, nil, repo.SaveAttachment(ctx, att))

	loaded, err := repo.GenAttachment(ctx, "gen-uuid")
	assert.Equal(t, nil, err)
	assert.Equal(t, "gen-uuid", loaded.AttachmentId)
	assert.Equal(t, "gen.pdf", loaded.FineName)
	assert.Equal(t, "application/pdf", loaded.ContentType)
	assert.Equal(t, []byte("gen-bytes"), loaded.Content)
	assert.Equal(t, testBucketName, loaded.Metadata[attachment.MetadataKeyBucket])
	assert.Equal(t, "2024/04/12/budget-gen_EXPENSE/gen-uuid", loaded.Metadata[attachment.MetadataKeyObjectKey])
}

func TestDeleteAttachmentRemovesContentAndMetadata(t *testing.T) {
	repo := newCompositeRepository("delete-uuid")

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:    "budget-del",
			BudgetType:  "expense",
			Date:        testutils.SafeDateFor("10/07/2024"),
			Owner:       "del-user",
			FineName:    "del.pdf",
			ContentType: "application/pdf",
		},
		Content: []byte("delete-me"),
	}
	ctx := testutils.NewStubbedContextWith("del-user")
	assert.Equal(t, nil, repo.SaveAttachment(ctx, att))

	expectedKey := "2024/07/10/budget-del_EXPENSE/delete-uuid"

	err := repo.DeleteAttachment(ctx, "delete-uuid")
	assert.Equal(t, nil, err)

	_, err = s3Client.GetObject(context.Background(), &aws_s3.GetObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(expectedKey),
	})
	assert.NotEqual(t, nil, err)

	item, err := dynamoClient.GetItem(context.Background(), &aws_dynamodb.GetItemInput{
		TableName: aws.String(testTableName),
		Key: map[string]dynamotypes.AttributeValue{
			"pk":            &dynamotypes.AttributeValueMemberS{Value: "budget-del_EXPENSE"},
			"attachment_id": &dynamotypes.AttributeValueMemberS{Value: "delete-uuid"},
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, 0, len(item.Item))
}

func TestDeleteAttachmentReturnsErrorWhenNotOwner(t *testing.T) {
	repo := newCompositeRepository("acl-uuid")

	att := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			BudgetId:    "budget-acl-comp",
			BudgetType:  "revenue",
			Date:        testutils.SafeDateFor("11/07/2024"),
			Owner:       "owner-real",
			FineName:    "acl.pdf",
			ContentType: "application/pdf",
		},
		Content: []byte("hands off"),
	}
	assert.Equal(t, nil, repo.SaveAttachment(testutils.NewStubbedContextWith("owner-real"), att))

	err := repo.DeleteAttachment(testutils.NewStubbedContextWith("intruder"), "acl-uuid")
	assert.NotEqual(t, nil, err)

	expectedKey := "2024/07/11/budget-acl-comp_REVENUE/acl-uuid"
	_, err = s3Client.GetObject(context.Background(), &aws_s3.GetObjectInput{
		Bucket: aws.String(testBucketName),
		Key:    aws.String(expectedKey),
	})
	assert.Equal(t, nil, err)
}

func TestFindAllAttachmentReturnsOnlyOwnerItemsForPartition(t *testing.T) {
	repo := newCompositeRepository("findall-uuid")

	mine := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "mine-1",
			BudgetId:     "budget-findall",
			BudgetType:   "revenue",
			Date:         testutils.SafeDateFor("01/05/2024"),
			Owner:        "find-user",
			FineName:     "mine.pdf",
			ContentType:  "application/pdf",
		},
		Content: []byte("mine"),
	}
	other := &attachment.Attachment{
		AttachmentMetadata: attachment.AttachmentMetadata{
			AttachmentId: "other-1",
			BudgetId:     "budget-findall",
			BudgetType:   "revenue",
			Date:         testutils.SafeDateFor("02/05/2024"),
			Owner:        "other-user",
			FineName:     "other.pdf",
			ContentType:  "application/pdf",
		},
		Content: []byte("other"),
	}
	ctx := testutils.NewStubbedContextWith("find-user")
	assert.Equal(t, nil, repo.SaveAttachment(ctx, mine))
	assert.Equal(t, nil, repo.SaveAttachment(testutils.NewStubbedContextWith("other-user"), other))

	got, err := repo.FindAllAttachment(ctx, "budget-findall", attachment.Revenue)
	assert.Equal(t, nil, err)
	assert.Equal(t, 1, len(got))
	assert.Equal(t, "mine-1", got[0].AttachmentId)
	assert.Equal(t, "find-user", got[0].Owner)
	assert.Equal(t, testBucketName, got[0].Metadata[attachment.MetadataKeyBucket])
}

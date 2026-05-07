//go:build test

package attachment

import (
	"context"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/go-playground/assert/v2"
	attachmentdynamo "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/dynamodb"
	attachments3 "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/s3"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/internal/testutils"
)

const (
	testTableName  = "BUDGET_ATTACHMENT_METADATA_COMPOSITE_TABLE"
	testBucketName = "budget-attachment-composite-staging"
)

var dynamoClient, _ = newLocalDynamoDBClient()
var s3Client, _ = newLocalS3Client()

func newLocalDynamoDBClient() (*aws_dynamodb.Client, error) {
	cfg, err := aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		aws_config.WithRegion("eu-central-1"),
		aws_config.WithBaseEndpoint("http://localhost:4566"),
	)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return aws_dynamodb.NewFromConfig(cfg), err
}

func newLocalS3Client() (*aws_s3.Client, error) {
	cfg, err := aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		aws_config.WithRegion("eu-central-1"),
		aws_config.WithBaseEndpoint("http://localhost:4566"),
	)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return aws_s3.NewFromConfig(cfg, func(o *aws_s3.Options) {
		o.UsePathStyle = true
	}), err
}

func setupCompositeTable() error {
	teardownCompositeTable()
	_, err := dynamoClient.CreateTable(context.TODO(), &aws_dynamodb.CreateTableInput{
		TableName: aws.String(testTableName),
		AttributeDefinitions: []dynamotypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("range_key"), AttributeType: dynamotypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamotypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamotypes.KeyTypeHash},
			{AttributeName: aws.String("range_key"), KeyType: dynamotypes.KeyTypeRange},
		},
		BillingMode: dynamotypes.BillingModePayPerRequest,
	})
	if err != nil {
		var resourceInUse *dynamotypes.ResourceInUseException
		if !errors.As(err, &resourceInUse) {
			return err
		}
	}
	return nil
}

func teardownCompositeTable() error {
	_, err := dynamoClient.DeleteTable(context.TODO(), &aws_dynamodb.DeleteTableInput{
		TableName: aws.String(testTableName),
	})
	return err
}

func setupCompositeBucket() error {
	teardownCompositeBucket()
	_, err := s3Client.CreateBucket(context.TODO(), &aws_s3.CreateBucketInput{
		Bucket: aws.String(testBucketName),
		CreateBucketConfiguration: &s3types.CreateBucketConfiguration{
			LocationConstraint: s3types.BucketLocationConstraintEuCentral1,
		},
	})
	if err != nil {
		var alreadyOwned *s3types.BucketAlreadyOwnedByYou
		var alreadyExists *s3types.BucketAlreadyExists
		if errors.As(err, &alreadyOwned) || errors.As(err, &alreadyExists) {
			return nil
		}
		return err
	}
	return nil
}

func teardownCompositeBucket() error {
	out, err := s3Client.ListObjectsV2(context.TODO(), &aws_s3.ListObjectsV2Input{
		Bucket: aws.String(testBucketName),
	})
	if err == nil {
		for _, obj := range out.Contents {
			_, _ = s3Client.DeleteObject(context.TODO(), &aws_s3.DeleteObjectInput{
				Bucket: aws.String(testBucketName),
				Key:    obj.Key,
			})
		}
	}
	_, err = s3Client.DeleteBucket(context.TODO(), &aws_s3.DeleteBucketInput{
		Bucket: aws.String(testBucketName),
	})
	return err
}

func TestMain(m *testing.M) {
	setupCompositeTable()
	setupCompositeBucket()
	code := m.Run()
	teardownCompositeBucket()
	teardownCompositeTable()
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
			"pk":        &dynamotypes.AttributeValueMemberS{Value: "budget-abc_EXPENSE"},
			"range_key": &dynamotypes.AttributeValueMemberS{Value: "generated-uuid-1"},
		},
	})
	assert.Equal(t, nil, err)
	expectedFileLocation := testBucketName + "/" + expectedKey
	assert.Equal(t, expectedFileLocation, itemOut.Item["file_location"].(*dynamotypes.AttributeValueMemberS).Value)
	assert.Equal(t, "generated-uuid-1", itemOut.Item["attachment_id"].(*dynamotypes.AttributeValueMemberS).Value)
	assert.Equal(t, "testuser", itemOut.Item["owner"].(*dynamotypes.AttributeValueMemberS).Value)
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
			"pk":        &dynamotypes.AttributeValueMemberS{Value: "budget-xyz_REVENUE"},
			"range_key": &dynamotypes.AttributeValueMemberS{Value: "preexisting-id"},
		},
	})
	assert.Equal(t, nil, err)
	assert.Equal(t, "new.pdf", itemOut.Item["file_name"].(*dynamotypes.AttributeValueMemberS).Value)
}

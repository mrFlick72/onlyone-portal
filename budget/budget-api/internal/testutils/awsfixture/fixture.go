//go:build test

// Package awsfixture provides LocalStack-targeted setup/teardown helpers
// shared by attachment adapter tests (DynamoDB metadata table + S3 bucket).
package awsfixture

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	dynamotypes "github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const localStackEndpoint = "http://localhost:4566"

func localStackConfig() (aws.Config, error) {
	return aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		aws_config.WithRegion("eu-central-1"),
		aws_config.WithBaseEndpoint(localStackEndpoint),
	)
}

func NewLocalDynamoDBClient() *dynamodb.Client {
	cfg, err := localStackConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return dynamodb.NewFromConfig(cfg)
}

func NewLocalS3Client() *aws_s3.Client {
	cfg, err := localStackConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return aws_s3.NewFromConfig(cfg, func(o *aws_s3.Options) {
		o.UsePathStyle = true
	})
}

// SetupAttachmentTable creates the DynamoDB table used by attachment metadata
// repositories: pk(HASH)+attachment_id(RANGE) with a GSI on attachment_id.
// An already-existing table is treated as success.
func SetupAttachmentTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.CreateTable(ctx, &dynamodb.CreateTableInput{
		TableName: aws.String(tableName),
		AttributeDefinitions: []dynamotypes.AttributeDefinition{
			{AttributeName: aws.String("pk"), AttributeType: dynamotypes.ScalarAttributeTypeS},
			{AttributeName: aws.String("attachment_id"), AttributeType: dynamotypes.ScalarAttributeTypeS},
		},
		KeySchema: []dynamotypes.KeySchemaElement{
			{AttributeName: aws.String("pk"), KeyType: dynamotypes.KeyTypeHash},
			{AttributeName: aws.String("attachment_id"), KeyType: dynamotypes.KeyTypeRange},
		},
		GlobalSecondaryIndexes: []dynamotypes.GlobalSecondaryIndex{
			{
				IndexName: aws.String(tableName + "_GLOBAL_INDEX"),
				KeySchema: []dynamotypes.KeySchemaElement{
					{AttributeName: aws.String("attachment_id"), KeyType: dynamotypes.KeyTypeHash},
				},
				Projection: &dynamotypes.Projection{
					ProjectionType: dynamotypes.ProjectionTypeAll,
				},
			},
		},
		BillingMode: dynamotypes.BillingModePayPerRequest,
	})
	if err != nil {
		var resourceInUse *dynamotypes.ResourceInUseException
		if errors.As(err, &resourceInUse) {
			return nil
		}
		return err
	}
	return nil
}

func TeardownTable(ctx context.Context, client *dynamodb.Client, tableName string) error {
	_, err := client.DeleteTable(ctx, &dynamodb.DeleteTableInput{
		TableName: aws.String(tableName),
	})
	return err
}

// SetupBucket creates an S3 bucket with eu-central-1 location constraint.
// An already-existing bucket (owned or otherwise) is treated as success.
func SetupBucket(ctx context.Context, client *aws_s3.Client, bucketName string) error {
	_, err := client.CreateBucket(ctx, &aws_s3.CreateBucketInput{
		Bucket: aws.String(bucketName),
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

// TeardownBucket drains all objects then deletes the bucket. List/object-delete
// errors are swallowed so a missing bucket still results in a clean exit code.
func TeardownBucket(ctx context.Context, client *aws_s3.Client, bucketName string) error {
	out, err := client.ListObjectsV2(ctx, &aws_s3.ListObjectsV2Input{
		Bucket: aws.String(bucketName),
	})
	if err == nil {
		for _, obj := range out.Contents {
			_, _ = client.DeleteObject(ctx, &aws_s3.DeleteObjectInput{
				Bucket: aws.String(bucketName),
				Key:    obj.Key,
			})
		}
	}
	_, err = client.DeleteBucket(ctx, &aws_s3.DeleteBucketInput{
		Bucket: aws.String(bucketName),
	})
	return err
}

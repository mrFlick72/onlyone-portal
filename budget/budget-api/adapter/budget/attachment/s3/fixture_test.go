//go:build test

package s3

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

const TestBucket = "budget-attachment-staging"

var s3Client, _ = newS3Client()

func newS3Client() (*aws_s3.Client, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		config.WithRegion("eu-central-1"),
		config.WithBaseEndpoint("http://localhost:4566"),
	)
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return aws_s3.NewFromConfig(cfg, func(o *aws_s3.Options) {
		o.UsePathStyle = true
	}), err
}

func setupTestBucket() error {
	teardownTestBucket()
	_, err := s3Client.CreateBucket(context.TODO(), &aws_s3.CreateBucketInput{
		Bucket: aws.String(TestBucket),
		CreateBucketConfiguration: &types.CreateBucketConfiguration{
			LocationConstraint: types.BucketLocationConstraintEuCentral1,
		},
	})
	if err != nil {
		var alreadyOwned *types.BucketAlreadyOwnedByYou
		var alreadyExists *types.BucketAlreadyExists
		if errors.As(err, &alreadyOwned) || errors.As(err, &alreadyExists) {
			return nil
		}
		return err
	}
	return nil
}

func teardownTestBucket() error {
	out, err := s3Client.ListObjectsV2(context.TODO(), &aws_s3.ListObjectsV2Input{
		Bucket: aws.String(TestBucket),
	})
	if err == nil {
		for _, obj := range out.Contents {
			_, _ = s3Client.DeleteObject(context.TODO(), &aws_s3.DeleteObjectInput{
				Bucket: aws.String(TestBucket),
				Key:    obj.Key,
			})
		}
	}
	_, err = s3Client.DeleteBucket(context.TODO(), &aws_s3.DeleteBucketInput{
		Bucket: aws.String(TestBucket),
	})
	return err
}

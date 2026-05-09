//go:build test

package s3

import (
	"context"

	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/mrflick72/budget/budget-api/internal/testutils/awsfixture"
)

const TestBucket = "budget-attachment-staging"

var s3Client *aws_s3.Client = awsfixture.NewLocalS3Client()

func setupTestBucket() error {
	_ = awsfixture.TeardownBucket(context.TODO(), s3Client, TestBucket)
	return awsfixture.SetupBucket(context.TODO(), s3Client, TestBucket)
}

func teardownTestBucket() error {
	return awsfixture.TeardownBucket(context.TODO(), s3Client, TestBucket)
}

//go:build test

// Package awsfixture provides LocalStack-targeted setup/teardown helpers
// shared by attachment adapter tests (DynamoDB metadata table + S3 bucket).
package cypto

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

const localStackEndpoint = "http://localhost:4566"

func localStackConfig() (aws.Config, error) {
	return aws_config.LoadDefaultConfig(context.TODO(),
		aws_config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		aws_config.WithRegion("eu-central-1"),
		aws_config.WithBaseEndpoint(localStackEndpoint),
	)
}

func newLocalKmsClient() *kms.Client {
	cfg, err := localStackConfig()
	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return kms.NewFromConfig(cfg)
}

package awsclient

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-sdk-go-v2/otelaws"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
)

// LoadDefaultConfig wraps aws_config.LoadDefaultConfig and appends OTel
// middleware to every AWS API call when otel.enabled=true.
// When otel.enabled=false a plain aws.Config is returned with no overhead.
func LoadDefaultConfig(ctx context.Context, optFns ...func(*aws_config.LoadOptions) error) (aws.Config, error) {
	cfg, err := aws_config.LoadDefaultConfig(ctx, optFns...)
	if err != nil {
		return cfg, err
	}
	if config.GetConfigurationManagerInstance().GetConfigBoolFor("otel.enabled") {
		otelaws.AppendMiddlewares(&cfg.APIOptions)
	}
	return cfg, nil
}

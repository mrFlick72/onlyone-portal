package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploggrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/log/global"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials/insecure"
)

func setupLoggerProvider(ctx context.Context, cfg otelConfig, res *resource.Resource) (ShutdownFunc, error) {
	logger.LogInfofFor("Setting up OTel LoggerProvider: service=%s protocol=%s endpoint=%s", cfg.ServiceName, cfg.Protocol, cfg.Endpoint)

	exp, err := buildLogExporter(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("otel: log exporter: %w", err)
	}

	lp := sdklog.NewLoggerProvider(
		sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
		sdklog.WithResource(res),
	)
	global.SetLoggerProvider(lp)

	return func(ctx context.Context) error {
		logger.LogInfofFor("Shutting down OTel LoggerProvider")
		return lp.Shutdown(ctx)
	}, nil
}

func buildLogExporter(ctx context.Context, cfg otelConfig) (sdklog.Exporter, error) {
	switch cfg.Protocol {
	case "grpc":
		opts := []otlploggrpc.Option{otlploggrpc.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlploggrpc.WithTLSCredentials(insecure.NewCredentials()))
		}
		return otlploggrpc.New(ctx, opts...)
	case "http", "":
		opts := []otlploghttp.Option{otlploghttp.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlploghttp.WithInsecure())
		}
		return otlploghttp.New(ctx, opts...)
	default:
		return nil, fmt.Errorf("unsupported protocol %q (must be \"http\" or \"grpc\")", cfg.Protocol)
	}
}

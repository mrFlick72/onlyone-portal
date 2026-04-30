package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"google.golang.org/grpc/credentials/insecure"
)

func setupMeterProvider(ctx context.Context, cfg otelConfig, res *resource.Resource) (ShutdownFunc, error) {
	logger.LogInfofFor("Setting up OTel MeterProvider: service=%s protocol=%s endpoint=%s", cfg.ServiceName, cfg.Protocol, cfg.Metrics.Endpoint)

	exp, err := buildMetricExporter(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("otel: metric exporter: %w", err)
	}

	mp := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(exp)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(mp)

	return func(ctx context.Context) error {
		logger.LogInfofFor("Shutting down OTel MeterProvider")
		return mp.Shutdown(ctx)
	}, nil
}

func buildMetricExporter(ctx context.Context, cfg otelConfig) (metric.Exporter, error) {
	switch cfg.Protocol {
	case "grpc":
		opts := []otlpmetricgrpc.Option{otlpmetricgrpc.WithEndpoint(cfg.Metrics.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlpmetricgrpc.WithTLSCredentials(insecure.NewCredentials()))
		}
		return otlpmetricgrpc.New(ctx, opts...)
	case "http", "":
		opts := []otlpmetrichttp.Option{otlpmetrichttp.WithEndpoint(cfg.Metrics.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlpmetrichttp.WithInsecure())
		}
		return otlpmetrichttp.New(ctx, opts...)
	default:
		return nil, fmt.Errorf("unsupported protocol %q (must be \"http\" or \"grpc\")", cfg.Protocol)
	}
}

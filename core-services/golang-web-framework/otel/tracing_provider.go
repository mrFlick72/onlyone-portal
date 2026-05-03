package otel

import (
	"context"
	"fmt"

	gootel "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials/insecure"
)

// SetupTracerProvider initialises only the TracerProvider.
// Prefer Setup() for full OTel initialisation.
func SetupTracerProvider(ctx context.Context) (ShutdownFunc, error) {
	cfg := loadOtelConfig()
	if !cfg.Enabled {
		logger.LogInfofFor("OTel tracing disabled")
		return noopShutdown, nil
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown-service"
	}
	res, err := buildResource(cfg)
	if err != nil {
		return nil, err
	}
	return setupTraceProvider(ctx, cfg, res)
}

func setupTraceProvider(ctx context.Context, cfg otelConfig, res *resource.Resource) (ShutdownFunc, error) {
	logger.LogInfofFor("Setting up OTel TracerProvider: service=%s protocol=%s endpoint=%s", cfg.ServiceName, cfg.Protocol, cfg.Endpoint)

	exp, err := buildTracExporter(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("otel: span exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
	gootel.SetTracerProvider(tp)
	gootel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		logger.LogInfofFor("Shutting down OTel TracerProvider")
		return tp.Shutdown(ctx)
	}, nil
}

func buildTracExporter(ctx context.Context, cfg otelConfig) (sdktrace.SpanExporter, error) {
	switch cfg.Protocol {
	case "grpc":
		opts := []otlptracegrpc.Option{otlptracegrpc.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlptracegrpc.WithTLSCredentials(insecure.NewCredentials()))
		}
		return otlptracegrpc.New(ctx, opts...)
	case "http", "":
		opts := []otlptracehttp.Option{otlptracehttp.WithEndpoint(cfg.Endpoint)}
		if cfg.Insecure {
			opts = append(opts, otlptracehttp.WithInsecure())
		}
		return otlptracehttp.New(ctx, opts...)
	default:
		return nil, fmt.Errorf("unsupported protocol %q (must be \"http\" or \"grpc\")", cfg.Protocol)
	}
}

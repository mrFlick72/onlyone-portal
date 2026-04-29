package otel

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"google.golang.org/grpc/credentials/insecure"
)

// ShutdownFunc flushes and closes OTel providers. Must be called on process exit.
type ShutdownFunc func(ctx context.Context) error

func noopShutdown(_ context.Context) error { return nil }

// Setup is the unified entry point. It initialises the global TracerProvider,
// MeterProvider, and LoggerProvider from YAML config.
// When otel.enabled=false, installs nothing and returns a no-op shutdown.
func Setup(ctx context.Context) (ShutdownFunc, error) {
	cfg := loadOtelConfig()
	if !cfg.Enabled {
		tracingLogger.LogInfofFor("OTel disabled")
		return noopShutdown, nil
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown-service"
	}

	res, err := buildResource(cfg)
	if err != nil {
		return nil, err
	}

	var shutdowns []ShutdownFunc

	traceShutdown, err := setupTraceProvider(ctx, cfg, res)
	if err != nil {
		return nil, err
	}
	shutdowns = append(shutdowns, traceShutdown)

	metricsShutdown, err := setupMeterProvider(ctx, cfg, res)
	if err != nil {
		return nil, err
	}
	shutdowns = append(shutdowns, metricsShutdown)

	logShutdown, err := setupLoggerProvider(ctx, cfg, res)
	if err != nil {
		return nil, err
	}
	shutdowns = append(shutdowns, logShutdown)

	return combineShutdowns(shutdowns), nil
}

// SetupTracerProvider initialises only the TracerProvider.
// Prefer Setup() for full OTel initialisation.
func SetupTracerProvider(ctx context.Context) (ShutdownFunc, error) {
	cfg := loadOtelConfig()
	if !cfg.Enabled {
		tracingLogger.LogInfofFor("OTel tracing disabled")
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

func buildResource(cfg otelConfig) (*resource.Resource, error) {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(semconv.SchemaURL, semconv.ServiceName(cfg.ServiceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("tracing: resource: %w", err)
	}
	return res, nil
}

func combineShutdowns(shutdowns []ShutdownFunc) ShutdownFunc {
	return func(ctx context.Context) error {
		var errs []error
		for _, shutdown := range shutdowns {
			if err := shutdown(ctx); err != nil {
				errs = append(errs, err)
			}
		}
		return errors.Join(errs...)
	}
}

func setupTraceProvider(ctx context.Context, cfg otelConfig, res *resource.Resource) (ShutdownFunc, error) {
	tracingLogger.LogInfofFor("Setting up OTel TracerProvider: service=%s protocol=%s endpoint=%s", cfg.ServiceName, cfg.Protocol, cfg.Endpoint)

	exp, err := buildExporter(ctx, cfg)
	if err != nil {
		return nil, fmt.Errorf("tracing: span exporter: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.AlwaysSample())),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return func(ctx context.Context) error {
		tracingLogger.LogInfofFor("Shutting down OTel TracerProvider")
		return tp.Shutdown(ctx)
	}, nil
}

// buildExporter builds a span exporter (package-private for test access).
func buildExporter(ctx context.Context, cfg otelConfig) (sdktrace.SpanExporter, error) {
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

package otel

import (
	"context"
	"errors"
	"fmt"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
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
		logger.LogInfofFor("OTel disabled")
		return noopShutdown, nil
	}
	if cfg.ServiceName == "" {
		cfg.ServiceName = "unknown-service"
	}

	res, err := buildResource(cfg)
	if err != nil {
		logger.LogErrorfFor("Error in buildResource: %v", err)
		return nil, err
	}

	var shutdowns []ShutdownFunc

	traceShutdown, err := setupTraceProvider(ctx, cfg, res)
	if err != nil {
		logger.LogErrorfFor("Error in setupTraceProvider: %v", err)
		return nil, err
	}
	shutdowns = append(shutdowns, traceShutdown)

	metricsShutdown, err := setupMeterProvider(ctx, cfg, res)
	if err != nil {
		logger.LogErrorfFor("Error in setupMeterProvider: %v", err)
		return nil, err
	}
	shutdowns = append(shutdowns, metricsShutdown)

	logShutdown, err := setupLoggerProvider(ctx, cfg, res)
	if err != nil {
		logger.LogErrorfFor("Error setupLoggerProvider: %v", err)
		return nil, err
	}
	shutdowns = append(shutdowns, logShutdown)

	return combineShutdowns(shutdowns), nil
}

func buildResource(cfg otelConfig) (*resource.Resource, error) {
	// Empty schema URL lets Merge inherit the URL from resource.Default() without
	// conflicting with whichever semconv version the SDK ships with.
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes("", attribute.String("service.name", cfg.ServiceName)),
	)
	if err != nil {
		return nil, fmt.Errorf("otel: resource: %w", err)
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

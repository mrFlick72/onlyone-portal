package otel

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/sdk/resource"
)

// partialSetupRollbackTimeout bounds the best-effort cleanup that runs when
// Setup fails after one or more providers have already been installed
// globally. Kept short because we're aborting boot — we want to stop the
// already-running goroutines and let the caller proceed to bail out, not
// wait minutes for an exporter that's clearly unreachable.
const partialSetupRollbackTimeout = 5 * time.Second

// ShutdownFunc flushes and closes OTel providers. Must be called on process exit.
type ShutdownFunc func(ctx context.Context) error

func noopShutdown(_ context.Context) error { return nil }

// Setup is the unified entry point. It initialises the global TracerProvider,
// MeterProvider, and LoggerProvider from YAML config.
// When otel.enabled=false, installs nothing and returns a no-op shutdown.
//
// Each setup* call installs its provider globally as a side effect. If a
// later setup fails, we must shut down the already-installed ones — otherwise
// their background goroutines (BatchSpanProcessor, PeriodicReader, ...) leak
// for the lifetime of the process. The committed/rollback pattern below
// guarantees that on any error path the partial state is rolled back before
// returning.
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
	committed := false
	defer func() {
		if committed {
			return
		}
		rollbackCtx, cancel := context.WithTimeout(context.Background(), partialSetupRollbackTimeout)
		defer cancel()
		if err := combineShutdowns(shutdowns)(rollbackCtx); err != nil {
			logger.LogErrorfFor("partial-setup rollback failed: %v", err)
		}
	}()

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

	committed = true
	return combineShutdowns(shutdowns), nil
}

// buildResource constructs the OTel Resource with service name, host, and SDK
// telemetry attributes. Using resource.New with explicit detectors avoids the
// schema URL conflict that resource.Merge(resource.Default(), ...) produces when
// the SDK's bundled semconv version differs from a pinned import.
func buildResource(cfg otelConfig) (*resource.Resource, error) {
	res, err := resource.New(context.Background(),
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(attribute.String("service.name", cfg.ServiceName)),
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

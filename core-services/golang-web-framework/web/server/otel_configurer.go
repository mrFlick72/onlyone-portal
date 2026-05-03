package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/otel"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

type OTelConfigurer struct {
	wsp       *WebServerProvisioner
	shutdown  otel.ShutdownFunc
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewOTelConfigurer(wsp *WebServerProvisioner) WebServerConfigurer {
	ctx, cancel := context.WithCancel(context.Background())
	configurer := &OTelConfigurer{
		wsp:       wsp,
		ctx:       ctx,
		cancelCtx: cancel,
	}
	wsp.configurers = append(wsp.configurers, configurer)

	return configurer
}

func (configurer *OTelConfigurer) Name() string {
	return "otel"
}

// Configure installs the global OTel providers and registers the otelgin
// middleware. Provider setup is non-fatal: on failure we log and fall back to
// a no-op shutdown so the service still boots without tracing.
func (configurer *OTelConfigurer) Configure() error {
	shutdown, err := otel.Setup(configurer.ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OTel setup failed (continuing without tracing): %v", err)
		shutdown = func(_ context.Context) error { return nil }
	}
	configurer.shutdown = shutdown

	// /management/* health probes are filtered out so liveness/readiness
	// pings don't pollute traces.
	serviceName := configurationManager.GetConfigFor("otel.service-name")
	configurer.wsp.engine.Use(otelgin.Middleware(serviceName,
		otelgin.WithFilter(func(req *http.Request) bool {
			return !strings.HasPrefix(req.URL.Path, "/management/")
		}),
	))

	return nil
}

// Dispose cancels the lifetime ctx (stopping background OTel goroutines) and
// flushes the providers within the caller-supplied deadline. The ctx parameter
// is the shutdown deadline, NOT configurer.ctx — the latter is the lifetime
// ctx passed to otel.Setup at boot, which by definition has no useful deadline
// for the flush.
func (configurer *OTelConfigurer) Dispose(ctx context.Context) error {
	defer configurer.cancelCtx()
	if configurer.shutdown != nil {
		return configurer.shutdown(ctx)
	}
	return nil
}

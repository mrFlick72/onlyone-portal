package server

import (
	"context"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/otel"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"strings"
)

type OtelWebServerConfigurer struct {
	wsp       *WebServerProvisioner
	shutdown  otel.ShutdownFunc
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewOtelWebServerConfigurer(wsp *WebServerProvisioner) WebServerConfigurer {
	ctx, cancel := context.WithCancel(context.Background())
	configurer := &OtelWebServerConfigurer{
		wsp:       wsp,
		ctx:       ctx,
		cancelCtx: cancel,
	}
	wsp.cancelContextFns = append(wsp.cancelContextFns, configurer)

	return configurer
}

func (configurer *OtelWebServerConfigurer) Name() string {
	return "otel"
}

// Configure installs the global OTel providers and registers the otelgin
// middleware. Provider setup is non-fatal: on failure we log and fall back to
// a no-op shutdown so the service still boots without tracing.
func (configurer *OtelWebServerConfigurer) Configure() error {
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

func (configurer *OtelWebServerConfigurer) Dispose() error {
	if configurer.shutdown != nil {
		if err := configurer.shutdown(configurer.ctx); err != nil {
			configurer.cancelCtx()
			return err
		}
	}
	configurer.cancelCtx()
	return nil
}

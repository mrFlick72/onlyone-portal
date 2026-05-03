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
	shutdown, err := otel.Setup(ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OTel setup failed (continuing without tracing): %v", err)
		shutdown = func(_ context.Context) error { return nil }
	}
	configurer := &OtelWebServerConfigurer{
		wsp:       wsp,
		shutdown:  shutdown,
		ctx:       ctx,
		cancelCtx: cancel,
	}
	wsp.cancelContextFns = append(wsp.cancelContextFns, configurer)

	return configurer
}

func (configurer *OtelWebServerConfigurer) Configure() error {
	serviceName := configurationManager.GetConfigFor("otel.service-name")
	// Lifecycle context for background goroutines (e.g. JWKS refresh).
	// Cancelled in Shutdown so they exit cleanly on process exit.

	configurer.wsp.engine.Use(otelgin.Middleware(serviceName,
		otelgin.WithFilter(func(req *http.Request) bool {
			return !strings.HasPrefix(req.URL.Path, "/management/")
		}),
	))

	return nil
}

func (configurer *OtelWebServerConfigurer) Dispose() error {
	err := configurer.shutdown(configurer.ctx)
	if err != nil {
		//todo add an error log message
		return err
	}
	configurer.cancelCtx()
	return nil
}

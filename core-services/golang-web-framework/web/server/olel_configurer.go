package server

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/otel"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"net/http"
	"strings"
)

type OtelWebServerConfigurer struct {
	engine    *gin.Engine
	shutdown  otel.ShutdownFunc
	cancelCtx context.CancelFunc
}

func NewOtelWebServerConfigurer(engine *gin.Engine) WebServerConfigurer {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown, err := otel.Setup(ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OTel setup failed (continuing without tracing): %v", err)
		shutdown = func(_ context.Context) error { return nil }
	}
	return &OtelWebServerConfigurer{
		engine:    engine,
		shutdown:  shutdown,
		cancelCtx: cancel,
	}
}

func (configurer *OtelWebServerConfigurer) Configure(ctx context.Context) (error, context.Context) {
	serviceName := configurationManager.GetConfigFor("otel.service-name")
	// Lifecycle context for background goroutines (e.g. JWKS refresh).
	// Cancelled in Shutdown so they exit cleanly on process exit.

	configurer.engine.Use(otelgin.Middleware(serviceName,
		otelgin.WithFilter(func(req *http.Request) bool {
			return !strings.HasPrefix(req.URL.Path, "/management/")
		}),
	))

	return nil, ctx
}

func (configurer *OtelWebServerConfigurer) Dispose(ctx context.Context) error {
	configurer.cancelCtx()
	return nil
}

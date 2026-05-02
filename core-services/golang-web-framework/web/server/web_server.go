package server

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/otel"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/magangement"
)

var configurationManager = config.GetConfigurationManagerInstance()
var web_server_logger = logging.GetLoggerInstanceForComponentByTypeName("WebServerProvisioner")

type WebServerProvisioner struct {
	router    *gin.Engine
	shutdown  otel.ShutdownFunc
	cancelCtx context.CancelFunc
}

func (wsp *WebServerProvisioner) ConfigureEngine() *gin.Engine {
	if wsp.router != nil {
		return wsp.router
	}

	// Lifecycle context for background goroutines (e.g. JWKS refresh).
	// Cancelled in Shutdown so they exit cleanly on process exit.
	ctx, cancel := context.WithCancel(context.Background())
	wsp.cancelCtx = cancel

	shutdown, err := otel.Setup(ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OTel setup failed (continuing without tracing): %v", err)
		shutdown = func(_ context.Context) error { return nil }
	}
	wsp.shutdown = shutdown

	router := gin.New()

	// 1. OTel: root server span wraps all subsequent middleware + handlers.
	//    /management/* health probes are filtered to avoid polluting traces.
	serviceName := configurationManager.GetConfigFor("otel.service-name")
	router.Use(otelgin.Middleware(serviceName,
		otelgin.WithFilter(func(req *http.Request) bool {
			return !strings.HasPrefix(req.URL.Path, "/management/")
		}),
	))

	// 2. Standard Gin middleware (inside the trace span)
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// 3. CORS
	router.Use(cors.New(cors.Config{
		AllowOrigins:     strings.Split(configurationManager.GetConfigFor("cors.allowed.origins"), ","),
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           60 * time.Minute,
	}))

	// 4. JWT/OAuth2 (auth failures visible as span events).
	//    JWKS misconfiguration is fatal at boot — without it every request would
	//    either be rejected or, worse, silently pass without verification.
	web_server_logger.LogInfofFor("Setting up OAuth2 middleware")
	oauth2, err := security.SetUpOAuth2(ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OAuth2 setup failed: %v", err)
		cancel()
		panic(fmt.Errorf("oauth2 setup: %w", err))
	}
	router.Use(oauth2)

	magangement.RegisterEndpoints(router)

	wsp.router = router
	return router
}

// Shutdown cancels background goroutines (JWKS refresh) and flushes OTel spans.
// Call after StartEngine() returns on process exit.
func (wsp *WebServerProvisioner) Shutdown(ctx context.Context) error {
	if wsp.cancelCtx != nil {
		wsp.cancelCtx()
	}
	if wsp.shutdown != nil {
		return wsp.shutdown(ctx)
	}
	return nil
}

func (wsp *WebServerProvisioner) StartEngine() error {
	port := configurationManager.GetConfigFor("server.port")
	serverBinder := fmt.Sprintf("0.0.0.0:%s", port)
	defer wsp.Shutdown(context.Background())

	if err := wsp.router.Run(serverBinder); err != nil {
		web_server_logger.LogErrorfFor("failed to run server: %v", err)
	}
	return nil
}

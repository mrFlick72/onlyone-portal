package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"strings"
	"syscall"
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

// HTTP server timeout defaults. Tunable via config keys server.read-timeout,
// server.write-timeout, server.idle-timeout, server.read-header-timeout,
// server.shutdown-timeout (Go duration strings, e.g. "30s", "2m").
const (
	defaultReadHeaderTimeout = 5 * time.Second
	defaultReadTimeout       = 30 * time.Second
	defaultWriteTimeout      = 30 * time.Second
	defaultIdleTimeout       = 120 * time.Second
	defaultShutdownTimeout   = 10 * time.Second
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

// StartEngine listens on server.port and blocks until the process receives
// SIGINT/SIGTERM or the listener fails. On signal, in-flight requests are
// drained for up to server.shutdown-timeout (default 10s) before OTel and
// JWKS background goroutines are stopped.
func (wsp *WebServerProvisioner) StartEngine() error {
	defer wsp.Shutdown(context.Background())

	port := configurationManager.GetConfigFor("server.port")
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	srv := &http.Server{
		Addr:              addr,
		Handler:           wsp.router,
		ReadHeaderTimeout: configurationManager.GetConfigDurationFor("server.read-header-timeout", defaultReadHeaderTimeout),
		ReadTimeout:       configurationManager.GetConfigDurationFor("server.read-timeout", defaultReadTimeout),
		WriteTimeout:      configurationManager.GetConfigDurationFor("server.write-timeout", defaultWriteTimeout),
		IdleTimeout:       configurationManager.GetConfigDurationFor("server.idle-timeout", defaultIdleTimeout),
	}

	serverErr := make(chan error, 1)
	go func() {
		web_server_logger.LogInfofFor("listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()

	sigCtx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serverErr:
		if err != nil {
			web_server_logger.LogErrorfFor("server failed: %v", err)
			return fmt.Errorf("server: %w", err)
		}
		return nil
	case <-sigCtx.Done():
		web_server_logger.LogInfofFor("shutdown signal received, draining connections")
	}

	shutdownTimeout := configurationManager.GetConfigDurationFor("server.shutdown-timeout", defaultShutdownTimeout)
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		web_server_logger.LogErrorfFor("graceful shutdown failed: %v", err)
		return fmt.Errorf("server shutdown: %w", err)
	}
	return nil
}

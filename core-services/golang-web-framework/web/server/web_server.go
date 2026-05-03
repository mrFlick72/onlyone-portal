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
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
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
	engine          *gin.Engine
	serverContext   context.Context
	cancelContextFn context.CancelFunc
}

/*
This is one utility interface to allow the server to inject behaviour like OAuth2, OTel and so on
This special feature can not simply apply a middleware
*/
type WebServerConfigurer interface {
	Configure() error
	Dispose() error
}

func (wsp *WebServerProvisioner) ConfigureEngine() *gin.Engine {
	if wsp.engine != nil {
		return wsp.engine
	}

	serverContext, cancelContextFn := context.WithCancel(context.Background())
	engine := gin.New()
	wsp.engine = engine
	wsp.serverContext = serverContext
	wsp.cancelContextFn = cancelContextFn

	// 1. OTel: root server span wraps all subsequent middleware + handlers.
	otelConfigurer := NewOtelWebServerConfigurer(engine)
	err := otelConfigurer.Configure()
	if err != nil {
		// todo fire an event that trigger the server shutdown
	}
	//    /management/* health probes are filtered to avoid polluting traces.

	// 2. Standard Gin middleware (inside the trace span)
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

	// 3. CORS
	engine.Use(corsCofigurer())

	// 4. JWT/OAuth2 (auth failures visible as span events).
	//    JWKS misconfiguration is fatal at boot — without it every request would
	//    either be rejected or, worse, silently pass without verification.
	web_server_logger.LogInfofFor("Setting up OAuth2 middleware")
	oauth2WebServerConfigurer := NewOauth2WebServerConfigurer(engine)
	err = oauth2WebServerConfigurer.Configure()
	if err != nil {
		// todo fire an event that trigger the server shutdown
	}

	magangement.RegisterEndpoints(engine)

	return engine
}

func corsCofigurer() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(configurationManager.GetConfigFor("cors.allowed.origins"), ","),
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           60 * time.Minute,
	})
}

// Shutdown cancels background goroutines (JWKS refresh) and flushes OTel spans.
// Call after StartEngine() returns on process exit.
func (wsp *WebServerProvisioner) Shutdown(ctx context.Context) error {
	wsp.cancelContextFn()
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
		Handler:           wsp.engine,
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

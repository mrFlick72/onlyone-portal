package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os/signal"
	"syscall"
	"time"

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
	engine      *gin.Engine
	configurers []WebServerConfigurer
}

// WebServerConfigurer lets the provisioner inject cross-cutting concerns
// (OAuth2, OTel, standard middleware, ...) that need lifecycle hooks beyond
// what a plain gin.HandlerFunc can express. Each configurer is registered with
// the provisioner at construction time and disposed on Shutdown.
type WebServerConfigurer interface {
	Name() string
	Configure() error
	Dispose(ctx context.Context) error
}

func (wsp *WebServerProvisioner) ConfigureEngine() *gin.Engine {
	if wsp.engine != nil {
		return wsp.engine
	}

	engine := gin.New()
	wsp.engine = engine

	// Order matters: OTel first so its server span wraps every later middleware
	// + handler; standard middleware (Logger/Recovery/CORS) inside the span so
	// panics and access logs are observed; OAuth2 last so auth failures show up
	// as span events on the already-open server span.
	configurers := []WebServerConfigurer{
		NewOTelConfigurer(wsp),
		NewStandardMiddlewareConfigurer(wsp),
		NewOAuth2Configurer(wsp),
	}
	for _, configurer := range configurers {
		web_server_logger.LogInfofFor("configuring %s", configurer.Name())
		if err := configurer.Configure(); err != nil {
			web_server_logger.LogErrorfFor("%s configure failed: %v", configurer.Name(), err)
			ctx, cancel := newShutdownContext()
			wsp.Shutdown(ctx)
			cancel()
			panic(fmt.Errorf("%s configure: %w", configurer.Name(), err))
		}
	}

	magangement.RegisterEndpoints(engine)

	return engine
}

// newShutdownContext returns a deadline-bounded context for any path that
// triggers Shutdown — graceful (SIGTERM), the ConfigureEngine panic path, and
// the StartEngine safety-net defer. Centralising it keeps server.shutdown-timeout
// the single source of truth for the shutdown budget.
func newShutdownContext() (context.Context, context.CancelFunc) {
	timeout := configurationManager.GetConfigDurationFor("server.shutdown-timeout", defaultShutdownTimeout)
	return context.WithTimeout(context.Background(), timeout)
}

// Shutdown disposes every registered configurer (cancelling background
// goroutines like JWKS refresh and flushing OTel spans) and resets the
// provisioner so it can be reconfigured. The ctx bounds disposal — in
// particular it caps how long the OTel exporter is allowed to flush before
// the process exits. A second call is a no-op because the slice has been
// emptied. Returns the joined errors of every Dispose that failed, so the
// caller (and the process exit code) can reflect a partial shutdown.
func (wsp *WebServerProvisioner) Shutdown(ctx context.Context) error {
	web_server_logger.LogInfofFor("Server Shutdown has been started ....")
	web_server_logger.LogInfofFor("There are %d configurer to be disposed", len(wsp.configurers))
	var errs []error
	for _, configurer := range wsp.configurers {
		if err := configurer.Dispose(ctx); err != nil {
			web_server_logger.LogErrorfFor("configurer dispose failed (continuing to clean up): %v", err)
			errs = append(errs, fmt.Errorf("%s dispose: %w", configurer.Name(), err))
		}
	}
	wsp.configurers = nil
	wsp.engine = nil
	web_server_logger.LogInfofFor("Server Shutdown completed ....")

	return errors.Join(errs...)
}

// StartEngine listens on server.port and blocks until the process receives
// SIGINT/SIGTERM or the listener fails. On signal, in-flight requests are
// drained for up to server.shutdown-timeout (default 10s) before OTel and
// JWKS background goroutines are stopped.
func (wsp *WebServerProvisioner) StartEngine() error {
	// Safety net: any early return (e.g. ListenAndServe failure) should still
	// dispose registered configurers within a bounded window. The graceful
	// path below clears the configurer slice, so on the happy path this defer
	// iterates an empty slice.
	defer func() {
		ctx, cancel := newShutdownContext()
		defer cancel()
		wsp.Shutdown(ctx)
	}()

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

	shutdownCtx, cancel := newShutdownContext()
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		web_server_logger.LogErrorfFor("graceful shutdown failed: %v", err)
		return fmt.Errorf("server shutdown: %w", err)
	}
	return wsp.Shutdown(shutdownCtx)
}

package server

import (
	"context"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

const healthCheckPath = "/management/health"

// StandardMiddlewareConfigurer registers the always-on Gin middleware:
// access logger, panic recovery, and CORS. It is registered between the OTel
// configurer and the OAuth2 configurer so the standard middleware runs inside
// the OTel server span (panics and access logs are observed) but before
// auth failures.
type StandardMiddlewareConfigurer struct {
	wsp *WebServerProvisioner
}

func NewStandardMiddlewareConfigurer(wsp *WebServerProvisioner) WebServerConfigurer {
	c := &StandardMiddlewareConfigurer{wsp: wsp}
	wsp.configurers = append(wsp.configurers, c)
	return c
}

func (c *StandardMiddlewareConfigurer) Name() string {
	return "standard-middleware"
}

func (c *StandardMiddlewareConfigurer) Configure() error {
	skipPaths := accessLogSkipPaths(configurationManager.GetConfigBoolFor("server.access-log.health-check-logging-enabled"))
	c.wsp.engine.Use(gin.LoggerWithConfig(gin.LoggerConfig{SkipPaths: skipPaths}))
	c.wsp.engine.Use(gin.Recovery())
	c.wsp.engine.Use(corsConfigurer())
	return nil
}

func (c *StandardMiddlewareConfigurer) Dispose(_ context.Context) error {
	return nil
}

// accessLogSkipPaths returns the gin access-log SkipPaths list. k8s liveness/
// readiness probes hit healthCheckPath every few seconds per pod, which
// drowns out real request activity in the access log — so it's skipped by
// default and only included when a service opts in via
// server.access-log.health-check-logging-enabled.
func accessLogSkipPaths(healthCheckLoggingEnabled bool) []string {
	if healthCheckLoggingEnabled {
		return nil
	}
	return []string{healthCheckPath}
}

func corsConfigurer() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(configurationManager.GetConfigFor("cors.allowed.origins"), ","),
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           60 * time.Minute,
	})
}

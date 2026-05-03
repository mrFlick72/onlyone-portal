package server

import (
	"context"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// StandardMiddlewareConfigurer registers the always-on Gin middleware:
// access logger, panic recovery, and CORS. It is registered between the OTel
// configurer and the OAuth2 configurer so the standard middleware runs inside
// the OTel server span (panics and access logs are observed) but before
// auth failures.
type StandardMiddlewareConfigurer struct {
	wsp *WebServerProvisioner
}

func NewStandardMiddlewareConfigurer(wsp *WebServerProvisioner) WebServerConfigurer {
	configurer := &StandardMiddlewareConfigurer{wsp: wsp}
	wsp.configurers = append(wsp.configurers, configurer)
	return configurer
}

func (configurer *StandardMiddlewareConfigurer) Name() string {
	return "standard-middleware"
}

func (configurer *StandardMiddlewareConfigurer) Configure() error {
	configurer.wsp.engine.Use(gin.Logger())
	configurer.wsp.engine.Use(gin.Recovery())
	configurer.wsp.engine.Use(corsConfigurer())
	return nil
}

func (configurer *StandardMiddlewareConfigurer) Dispose(_ context.Context) error {
	return nil
}

func corsConfigurer() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:     strings.Split(configurationManager.GetConfigFor("cors.allowed.origins"), ","),
		AllowMethods:     []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		MaxAge:           60 * time.Minute,
	})
}

package server

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/otel"
)

type OAuth2WebServerConfigurer struct {
	engine    *gin.Engine
	shutdown  otel.ShutdownFunc
	cancelCtx context.CancelFunc
}

func NewOauth2WebServerConfigurer(engine *gin.Engine) WebServerConfigurer {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown, err := otel.Setup(ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OTel setup failed (continuing without tracing): %v", err)
		shutdown = func(_ context.Context) error { return nil }
	}
	return &OAuth2WebServerConfigurer{
		engine:    engine,
		shutdown:  shutdown,
		cancelCtx: cancel,
	}
}

func (configurer *OAuth2WebServerConfigurer) Configure(ctx context.Context) (error, context.Context) {
	oauth2, err := security.SetUpOAuth2(ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OAuth2 setup failed: %v", err)
		panic(fmt.Errorf("oauth2 setup: %w", err))
	}
	configurer.engine.Use(oauth2)
	return nil, ctx
}

func (configurer *OAuth2WebServerConfigurer) Dispose(ctx context.Context) error {
	configurer.cancelCtx()
	return nil
}

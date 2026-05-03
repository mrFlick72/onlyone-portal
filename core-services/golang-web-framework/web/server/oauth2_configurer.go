package server

import (
	"context"
	"fmt"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/otel"
)

type OAuth2WebServerConfigurer struct {
	wsp       *WebServerProvisioner
	shutdown  otel.ShutdownFunc
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewOauth2WebServerConfigurer(wsp *WebServerProvisioner) WebServerConfigurer {
	ctx, cancel := context.WithCancel(context.Background())
	shutdown, err := otel.Setup(ctx)
	if err != nil {
		web_server_logger.LogErrorfFor("OTel setup failed (continuing without tracing): %v", err)
		shutdown = func(_ context.Context) error { return nil }
	}
	return &OAuth2WebServerConfigurer{
		wsp:       wsp,
		ctx:       ctx,
		shutdown:  shutdown,
		cancelCtx: cancel,
	}
}

func (configurer *OAuth2WebServerConfigurer) Configure() error {
	oauth2, err := security.SetUpOAuth2(context.Background())
	if err != nil {
		web_server_logger.LogErrorfFor("OAuth2 setup failed: %v", err)
		panic(fmt.Errorf("oauth2 setup: %w", err))
	}
	configurer.wsp.engine.Use(oauth2)
	return nil
}

func (configurer *OAuth2WebServerConfigurer) Dispose() error {
	err := configurer.shutdown(configurer.ctx)
	if err != nil {
		//todo add an error log message
		return err
	}
	configurer.cancelCtx()
	return nil
}

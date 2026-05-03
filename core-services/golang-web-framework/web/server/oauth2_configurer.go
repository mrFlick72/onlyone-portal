package server

import (
	"context"
	"fmt"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type OAuth2WebServerConfigurer struct {
	wsp       *WebServerProvisioner
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewOauth2WebServerConfigurer(wsp *WebServerProvisioner) WebServerConfigurer {
	ctx, cancel := context.WithCancel(context.Background())
	configurer := &OAuth2WebServerConfigurer{
		wsp:       wsp,
		ctx:       ctx,
		cancelCtx: cancel,
	}
	wsp.cancelContextFns = append(wsp.cancelContextFns, configurer)
	return configurer
}

func (configurer *OAuth2WebServerConfigurer) Name() string {
	return "oauth2"
}

func (configurer *OAuth2WebServerConfigurer) Configure() error {
	oauth2, err := security.SetUpOAuth2(configurer.ctx)
	if err != nil {
		return fmt.Errorf("oauth2 setup: %w", err)
	}
	configurer.wsp.engine.Use(oauth2)
	return nil
}

func (configurer *OAuth2WebServerConfigurer) Dispose() error {
	configurer.cancelCtx()
	return nil
}

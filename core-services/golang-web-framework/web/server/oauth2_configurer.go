package server

import (
	"context"
	"fmt"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type OAuth2Configurer struct {
	wsp       *WebServerProvisioner
	ctx       context.Context
	cancelCtx context.CancelFunc
}

func NewOAuth2Configurer(wsp *WebServerProvisioner) WebServerConfigurer {
	ctx, cancel := context.WithCancel(context.Background())
	configurer := &OAuth2Configurer{
		wsp:       wsp,
		ctx:       ctx,
		cancelCtx: cancel,
	}
	wsp.configurers = append(wsp.configurers, configurer)
	return configurer
}

func (configurer *OAuth2Configurer) Name() string {
	return "oauth2"
}

func (configurer *OAuth2Configurer) Configure() error {
	oauth2, err := security.SetUpOAuth2(configurer.ctx)
	if err != nil {
		return fmt.Errorf("oauth2 setup: %w", err)
	}
	configurer.wsp.engine.Use(oauth2)
	return nil
}

func (configurer *OAuth2Configurer) Dispose() error {
	configurer.cancelCtx()
	return nil
}

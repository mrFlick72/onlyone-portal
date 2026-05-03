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
	c := &OAuth2Configurer{
		wsp:       wsp,
		ctx:       ctx,
		cancelCtx: cancel,
	}
	wsp.configurers = append(wsp.configurers, c)
	return c
}

func (c *OAuth2Configurer) Name() string {
	return "oauth2"
}

func (c *OAuth2Configurer) Configure() error {
	oauth2, err := security.SetUpOAuth2(c.ctx)
	if err != nil {
		return fmt.Errorf("oauth2 setup: %w", err)
	}
	c.wsp.engine.Use(oauth2)
	return nil
}

func (c *OAuth2Configurer) Dispose(_ context.Context) error {
	c.cancelCtx()
	return nil
}

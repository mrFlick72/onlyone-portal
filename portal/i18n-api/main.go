package main

import (
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/portal/i18n-api/config"
	web "github.com/mrflick72/onlyone-portal/portal/i18n-api/web/i18n"
)

func main() {
	provisioner := server.WebServerProvisioner{}
	engine := provisioner.ConfigureEngine()
	factory := &server.GinContextToPlainContextFactory{}

	web.RegisterEndpoints(engine, factory, config.NewBundleRepository(), config.NewLanguageResolver())

	provisioner.StartEngine()
}

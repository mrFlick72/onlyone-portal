package main

import (
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/configuration"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/src/plan"
	api "github.com/mrflick72/onlyone-portal/plan/plan-service/src/web"
)

func main() {
	engine := server.WebServerProvisioner{}
	ginEngine := engine.ConfigureEngine()
	factory := &server.GinContextToPlainContextFactory{}

	api.RegisterEndpoints(ginEngine, factory, plan.NewPostgresTodoRepository(configuration.NewPostgresDSN()))

	engine.StartEngine()
}

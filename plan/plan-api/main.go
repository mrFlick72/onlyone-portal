package main

import (
	"github.com/mrflick72/onlyone-portal/plan/plan-service/config"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/web/plan"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func main() {
	engine := server.WebServerProvisioner{}
	ginEngine := engine.ConfigureEngine()
	factory := &server.GinContextToPlainContextFactory{}

	plan.RegisterEndpoints(ginEngine, factory, config.NewPlanRepository(config.NewPostgresDSN()))

	engine.StartEngine()
}

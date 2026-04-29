package main

import (
	plandb "github.com/mrflick72/onlyone-portal/plan/plan-service/adapter/plan/db"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/config"
	planapi "github.com/mrflick72/onlyone-portal/plan/plan-service/web/plan"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func main() {
	engine := server.WebServerProvisioner{}
	ginEngine := engine.ConfigureEngine()
	factory := &server.GinContextToPlainContextFactory{}

	planapi.RegisterEndpoints(ginEngine, factory, plandb.NewPlanRepository(config.NewPostgresDSN()))

	engine.StartEngine()
}

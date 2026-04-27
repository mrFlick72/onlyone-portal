package main

import (
	tododb "github.com/mrflick72/onlyone-portal/plan/plan-service/adapter/todo/db"
	"github.com/mrflick72/onlyone-portal/plan/plan-service/config"
	todoapi "github.com/mrflick72/onlyone-portal/plan/plan-service/web/todo"

	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func main() {
	engine := server.WebServerProvisioner{}
	ginEngine := engine.ConfigureEngine()
	factory := &server.GinContextToPlainContextFactory{}

	todoapi.RegisterEndpoints(ginEngine, factory, tododb.NewTodoRepository(config.NewPostgresDSN()))

	engine.StartEngine()
}

package main

import (
	"github.com/mrflick72/budget/budget-api/config"
	"github.com/mrflick72/budget/budget-api/web"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func main() {

	// Entry point for the budget API v2 service
	// Create a Gin router with default middleware (logger and recovery)
	engine := server.WebServerProvisioner{}

	ginEngine := engine.ConfigureEngine()
	GinContextToPlainContextFactory := &server.GinContextToPlainContextFactory{}
	web.RegisterEndpoints(ginEngine, GinContextToPlainContextFactory, config.NewBudgetExpenseActionsFacade())

	engine.StartEngine()
}

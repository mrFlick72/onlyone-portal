package main

import (
	"github.com/mrflick72/budget/budget-api/config"
	"github.com/mrflick72/budget/budget-api/web/budget/attachment"
	"github.com/mrflick72/budget/budget-api/web/budget/expense"
	"github.com/mrflick72/budget/budget-api/web/budget/revenue"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func main() {

	// Entry point for the budget API v2 service
	// Create a Gin router with default middleware (logger and recovery)
	engine := server.WebServerProvisioner{}

	ginEngine := engine.ConfigureEngine()
	GinContextToPlainContextFactory := &server.GinContextToPlainContextFactory{}
	expenseFacade, stopReclassificationListener := config.NewBudgetExpenseActionsFacade()
	defer stopReclassificationListener()
	expense.RegisterExpenseEndpoints(ginEngine, GinContextToPlainContextFactory, expenseFacade)
	revenue.RegisterRevenueEndpoints(ginEngine, GinContextToPlainContextFactory, config.NewRevenueActionsFacade())
	attachment.RegisterAttachmentEndpoints(ginEngine, GinContextToPlainContextFactory, config.NewAttachmentActionsFacade())

	engine.StartEngine()
}

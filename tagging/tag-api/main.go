package main

import (
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/adapter/dynamodb"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/web/api"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/web/server"
)

func main() {
	// Create a Gin router with default middleware (logger and recovery)
	engine := server.WebServerProvisioner{}

	router := engine.ConfigureEngine()
	
	api.RegisterEndpoints(router, dynamodb.NewTagDynamoDBRepository())

	engine.StartEngine()
}

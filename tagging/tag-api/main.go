package main

import (
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/adapter/dynamodb"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/web/api"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
)

func main() {

	// Entry point for the budget API v2 service
	// Create a Gin router with default middleware (logger and recovery)
	engine := server.WebServerProvisioner{}

	ginEngine := engine.ConfigureEngine()
	GinContextToPlainContextFactory := &server.GinContextToPlainContextFactory{}
	repository := dynamodb.NewTagDynamoDBRepository()
	findAllTagsAction := &domain.FindAllTags{Repository: repository}
	api.RegisterEndpoints(ginEngine, GinContextToPlainContextFactory, repository, findAllTagsAction)

	engine.StartEngine()
}

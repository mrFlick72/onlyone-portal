package config

import (
	"context"
	"net/http"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/google/uuid"
	"github.com/mrflick72/budget/budget-api/adapter/budget/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/tags/rest"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

var configurationManager = config.GetConfigurationManagerInstance()
var logger = logging.GetLoggerInstance()

func NewSearchTagRepository() tags.SearchTagRepository {
	return &rest.RestSearchTagRepository{
		Client:  &http.Client{},
		BaseURL: configurationManager.GetConfigFor("tag-api.base-url"),
	}
}

func NewBudgetExpenseRepository() expense.BudgetExpenseRepository {
	cfg, err := aws_config.LoadDefaultConfig(
		context.TODO(),
		aws_config.WithRegion("eu-central-1"),
	)

	if err != nil {
		logger.LogErrorfFor("unable to load SDK config: %s", err.Error())
		panic("unable to load SDK config, " + err.Error())
	}

	return &dynamodb.DynamoDbBudgetExpenseRepository{
		TableName: configurationManager.GetConfigFor("budget-api.dynamo-db.budget-expense.table-name"),
		Client:    aws_dynamodb.NewFromConfig(cfg),
		BudgetExpenseIdProvider: &dynamodb.DynamoDbBudgetExpenseIdProvider{
			SaltGenerator: func() string { return uuid.New().String() },
		},
		SearchTagRepository: NewSearchTagRepository(),
	}
}

func NewBudgetExpenseActionsFacade() expense.BudgetExpenseActions {
	budgetExpenseRepository := NewBudgetExpenseRepository()
	searchTagRepository := NewSearchTagRepository()
	createBudgetExpense := &expense.CreateBudgetExpense{
		Repository: budgetExpenseRepository,
	}
	updateBudgetExpense := &expense.UpdateBudgetExpense{
		Repository: budgetExpenseRepository,
	}
	findSpentBudget := &expense.FindSpentBudget{
		BudgetExpenseRepository: budgetExpenseRepository,
		SearchTagRepository:     searchTagRepository,
	}
	deleteBudgetExpense := &expense.DeleteBudgetExpense{
		Repository: budgetExpenseRepository,
	}

	return &expense.BudgetExpenseActionsFacade{
		CreateBudgetExpenseAction: createBudgetExpense,
		UpdateBudgetExpenseAction: updateBudgetExpense,
		FindSpentBudgetAction:     findSpentBudget,
		DeleteBudgetExpenseAction: deleteBudgetExpense,
	}
}

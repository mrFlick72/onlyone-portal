package config

import (
	"context"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	attachmentadapter "github.com/mrflick72/budget/budget-api/adapter/budget/attachment"
	attachmentdynamo "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/dynamodb"
	attachments3 "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/s3"
	"github.com/mrflick72/budget/budget-api/adapter/budget/expense/dynamodb"
	revenuedynamo "github.com/mrflick72/budget/budget-api/adapter/budget/revenue/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/tags/rest"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/awsclient"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/httpclient"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
)

var configurationManager = config.GetConfigurationManagerInstance()
var logger = logging.GetLoggerInstance()

func NewSearchTagRepository() tags.SearchTagRepository {
	delegate := rest.NewRestSearchTagRepository(
		httpclient.NewHTTPClient(),
		configurationManager.GetConfigFor("tag-api.base-url"),
	)

	return rest.NewRistrettoCachedSearchTagRepository(delegate)
}

func NewBudgetExpenseRepository() expense.BudgetExpenseRepository {
	cfg, err := awsclient.LoadDefaultConfig(
		context.TODO(),
		aws_config.WithRegion("eu-central-1"),
	)

	if err != nil {
		logger.LogErrorfFor("unable to load SDK config: %s", err.Error())
		panic("unable to load SDK config, " + err.Error())
	}

	return dynamodb.NewDynamoDbBudgetExpenseRepository(
		configurationManager.GetConfigFor("budget-api.dynamo-db.budget-expense.table-name"),
		aws_dynamodb.NewFromConfig(cfg),
		&dynamodb.DynamoDbBudgetExpenseIdProvider{
			SaltGenerator: func() string { return uuid.New().String() },
		},
		NewSearchTagRepository(),
	)

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

func NewRevenueRepository() revenue.RevenueRepository {
	cfg, err := awsclient.LoadDefaultConfig(
		context.TODO(),
		aws_config.WithRegion("eu-central-1"),
	)

	if err != nil {
		logger.LogErrorfFor("unable to load SDK config: %s", err.Error())
		panic("unable to load SDK config, " + err.Error())
	}

	return revenuedynamo.NewDynamoDbRevenueRepository(
		configurationManager.GetConfigFor("budget-api.dynamo-db.revenue.table-name"),
		aws_dynamodb.NewFromConfig(cfg),
		&revenuedynamo.DynamoDbRevenueIdProvider{
			SaltGenerator: func() string { return uuid.New().String() },
		},
	)
}

func NewRevenueActionsFacade() revenue.RevenueActions {
	revenueRepository := NewRevenueRepository()
	return &revenue.RevenueActionsFacade{
		CreateRevenueAction: &revenue.CreateRevenue{Repository: revenueRepository},
		UpdateRevenueAction: &revenue.UpdateRevenue{Repository: revenueRepository},
		FindRevenueAction:   &revenue.FindRevenue{Repository: revenueRepository},
		DeleteRevenueAction: &revenue.DeleteRevenue{Repository: revenueRepository},
	}
}

func NewAttachmentRepository() attachment.AttachmentRepository {
	cfg, err := awsclient.LoadDefaultConfig(
		context.TODO(),
		aws_config.WithRegion("eu-central-1"),
	)
	if err != nil {
		logger.LogErrorfFor("unable to load SDK config: %s", err.Error())
		panic("unable to load SDK config, " + err.Error())
	}

	idProvider := &attachmentdynamo.DynamoDbAttachmentIdProvider{
		UuidGenerator: func() string { return uuid.New().String() },
	}
	metadataRepository := attachmentdynamo.NewDynamoDbAttachmentMetadataRepository(
		configurationManager.GetConfigFor("budget-api.dynamo-db.attachment-metadata.table-name"),
		aws_dynamodb.NewFromConfig(cfg),
		idProvider,
	)
	contentRepository := attachments3.NewS3AttachmentContentRepository(
		configurationManager.GetConfigFor("budget-api.s3.attachment.bucket-name"),
		aws_s3.NewFromConfig(cfg),
	)

	return attachmentadapter.NewAwsCompositeAttachmentRepository(idProvider, metadataRepository, contentRepository)
}

func NewAttachmentActionsFacade() attachment.AttachmentActions {
	repository := NewAttachmentRepository()
	return &attachment.AttachmentActionsFacade{
		SaveAttachmentAction:   attachment.NewSaveAttachment(repository),
		GetAttachmentAction:    &attachment.GetAttachment{},
		DeleteAttachmentAction: &attachment.DeleteAttachment{},
	}
}

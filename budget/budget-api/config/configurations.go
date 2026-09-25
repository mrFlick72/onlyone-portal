package config

import (
	"context"
	"flag"
	"fmt"
	"strings"
	"time"

	aws_config "github.com/aws/aws-sdk-go-v2/config"
	aws_dynamodb "github.com/aws/aws-sdk-go-v2/service/dynamodb"
	aws_s3 "github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	attachmentadapter "github.com/mrflick72/budget/budget-api/adapter/budget/attachment"
	attachmentdynamo "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/dynamodb"
	attachments3 "github.com/mrflick72/budget/budget-api/adapter/budget/attachment/s3"
	"github.com/mrflick72/budget/budget-api/adapter/budget/expense/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/budget/expense/kafka"
	revenuedynamo "github.com/mrflick72/budget/budget-api/adapter/budget/revenue/dynamodb"
	scheduledexpensedynamo "github.com/mrflick72/budget/budget-api/adapter/budget/scheduledexpense/dynamodb"
	"github.com/mrflick72/budget/budget-api/adapter/budget/scheduledexpense/scheduler"
	"github.com/mrflick72/budget/budget-api/adapter/tags/rest"
	"github.com/mrflick72/budget/budget-api/domain/budget/attachment"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/awsclient"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/httpclient"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/web/server"
	"github.com/twmb/franz-go/pkg/kgo"
)

var configurationManager = config.GetConfigurationManagerInstance()
var logger = logging.GetLoggerInstance()

// newSearchTagRepository builds a tag-api-backed, Ristretto-cached tag
// repository scoped to a single tag Scope. The scope literal is supplied by the
// named wrappers below, keeping it defined once per scope at the wiring layer.
// Scope is an adapter/wiring concern and never reaches the domain. See
// docs/adr/0001-expense-scoped-tag-lookup-hardcoded-at-wiring.md and
// docs/adr/0002-revenue-tagging-mirrors-expense-without-events-or-totals.md.
func newSearchTagRepository(scope string) tags.SearchTagRepository {
	delegate := rest.NewRestSearchTagRepository(
		httpclient.NewHTTPClient(),
		configurationManager.GetConfigFor("tag-api.base-url"),
		scope,
	)

	return rest.NewRistrettoCachedSearchTagRepository(delegate, scope)
}

// NewExpenseSearchTagRepository resolves expense tags (GET /api/tags/scope/expense).
func NewExpenseSearchTagRepository() tags.SearchTagRepository {
	return newSearchTagRepository("expense")
}

// NewRevenueSearchTagRepository resolves revenue tags (GET /api/tags/scope/revenue).
func NewRevenueSearchTagRepository() tags.SearchTagRepository {
	return newSearchTagRepository("revenue")
}

func NewBudgetExpenseRepository(eventBus *expense.InternalEventBus) expense.BudgetExpenseRepository {
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
		NewExpenseSearchTagRepository(),
		eventBus,
		configurationManager.GetConfigBoolFor("budget-api.reclassification.enabled"),
	)

}

func NewNewKafkaBudgetExpenseEventPublisher() expense.BudgetExpenseEventPublisher {
	topic := configurationManager.GetConfigFor("budget-api.events.kafka.topic-name")
	brokerList := configurationManager.GetConfigFor("budget-api.events.kafka.brokers")
	batchMaxBytes := flag.Int("batch-max-bytes", 1000000, "the maximum batch size to allow per-partition (must be less than Kafka's max.message.bytes, producing)")
	recordBytes := flag.Int("record-bytes", 100, "bytes per record value (producing)")

	opts := []kgo.Opt{
		kgo.SeedBrokers(strings.Split(brokerList, ",")...),
		kgo.MaxBufferedRecords(250<<20 / *recordBytes + 1),
		kgo.MaxConcurrentFetches(3),
		// We have good compression, so we want to limit what we read
		// back because snappy deflation will balloon our memory usage.
		kgo.FetchMaxBytes(5 << 20),
		kgo.ProducerBatchMaxBytes(int32(*batchMaxBytes)),
	}
	client, err := kgo.NewClient(opts...)
	if err != nil || topic == "" || brokerList == "" {
		panic(fmt.Sprintf("Error during Kafka client configuration: %v", err))
	}
	return kafka.NewKafkaBudgetExpenseEventPublisher(topic, client)
}

// NewBudgetExpenseActionsFacade wires the expense actions and starts the
// reclassification listener. The returned stop function stops that listener and
// must be deferred by the caller so the goroutine is torn down with the process
// (after the HTTP server has drained, so no in-flight read can still publish).
// NewBudgetExpenseActionsFacade returns the concrete facade (it satisfies
// expense.BudgetExpenseActions for the endpoints) so main.go can also hand its
// CreateBudgetExpenseAction to the Scheduled Expense job.
func NewBudgetExpenseActionsFacade() (*expense.BudgetExpenseActionsFacade, func()) {
	eventBus := expense.NewEventBus()
	budgetExpenseRepository := NewBudgetExpenseRepository(eventBus)
	searchTagRepository := NewExpenseSearchTagRepository()
	eventPublisher := NewNewKafkaBudgetExpenseEventPublisher()
	createBudgetExpense := &expense.CreateBudgetExpense{
		Repository:     budgetExpenseRepository,
		EventPublisher: eventPublisher,
		Logger:         logging.GetLoggerInstanceForComponentByTypeName("CreateBudgetExpense"),
	}
	updateBudgetExpense := &expense.UpdateBudgetExpense{
		Repository:     budgetExpenseRepository,
		EventPublisher: eventPublisher,
		EventBus:       eventBus,
		Logger:         logging.GetLoggerInstanceForComponentByTypeName("UpdateBudgetExpense"),
	}
	go updateBudgetExpense.Listen()

	findSpentBudget := &expense.FindSpentBudget{
		BudgetExpenseRepository: budgetExpenseRepository,
		SearchTagRepository:     searchTagRepository,
		Logger:                  logging.GetLoggerInstanceForComponentByTypeName("FindSpentBudget"),
	}
	deleteBudgetExpense := &expense.DeleteBudgetExpense{
		Repository:     budgetExpenseRepository,
		EventPublisher: eventPublisher,
		Logger:         logging.GetLoggerInstanceForComponentByTypeName("DeleteBudgetExpense"),
	}

	facade := &expense.BudgetExpenseActionsFacade{
		CreateBudgetExpenseAction: createBudgetExpense,
		UpdateBudgetExpenseAction: updateBudgetExpense,
		FindSpentBudgetAction:     findSpentBudget,
		DeleteBudgetExpenseAction: deleteBudgetExpense,
	}
	return facade, eventBus.Close
}

func NewRevenueRepository() revenue.RevenueRepository {
	cfg, err := awsclient.LoadDefaultConfig(
		context.Background(),
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
		NewRevenueSearchTagRepository(),
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

func NewScheduledExpenseRepository() scheduledexpense.ScheduledExpenseRepository {
	cfg, err := awsclient.LoadDefaultConfig(
		context.Background(),
		aws_config.WithRegion("eu-central-1"),
	)

	if err != nil {
		logger.LogErrorfFor("unable to load SDK config: %s", err.Error())
		panic("unable to load SDK config, " + err.Error())
	}

	return scheduledexpensedynamo.NewDynamoDbScheduledExpenseRepository(
		configurationManager.GetConfigFor("budget-api.dynamo-db.scheduled-expense.table-name"),
		aws_dynamodb.NewFromConfig(cfg),
		&scheduledexpensedynamo.DynamoDbScheduledExpenseIdProvider{
			UuidGenerator: func() string { return uuid.New().String() },
		},
		// Scheduled Expense tags share expense's scope (see CONTEXT.md: "a tag
		// list (tag keys, as on BudgetExpense)") — not a separate scope.
		NewExpenseSearchTagRepository(),
	)
}

// NewScheduledExpenseJobConfigurer builds the in-process generation
// engine (ADR 0005) as a WebServerConfigurer for main.go to hand to
// WebServerProvisioner.RegisterConfigurer. createBudgetExpense must be the
// action held by the facade main.go already built — building another facade
// would start a second reclassification listener and Kafka client.
// Interval: budget-api.scheduled-expense.job.interval (Go duration,
// default 1h).
func NewScheduledExpenseJobConfigurer(createBudgetExpense *expense.CreateBudgetExpense) server.WebServerConfigurer {
	job := &scheduledexpense.ScheduledExpenseJob{
		Repository:          NewScheduledExpenseRepository(),
		CreateBudgetExpense: createBudgetExpense,
		Today:               date.Today,
		Logger:              logging.GetLoggerInstanceForComponentByTypeName("ScheduledExpenseJob"),
	}
	interval := configurationManager.GetConfigDurationFor("budget-api.scheduled-expense.job.interval", time.Hour)
	return scheduler.NewScheduledExpenseJobConfigurer(job, interval)
}

func NewScheduledExpenseActionsFacade() scheduledexpense.ScheduledExpenseActions {
	scheduledExpenseRepository := NewScheduledExpenseRepository()
	return &scheduledexpense.ScheduledExpenseActionsFacade{
		CreateScheduledExpenseAction: &scheduledexpense.CreateScheduledExpense{Repository: scheduledExpenseRepository},
		FindScheduledExpensesAction:  &scheduledexpense.FindScheduledExpenses{Repository: scheduledExpenseRepository},
		FindScheduledExpenseAction:   &scheduledexpense.FindScheduledExpense{Repository: scheduledExpenseRepository},
		UpdateScheduledExpenseAction: &scheduledexpense.UpdateScheduledExpense{Repository: scheduledExpenseRepository},
		DeleteScheduledExpenseAction: &scheduledexpense.DeleteScheduledExpense{Repository: scheduledExpenseRepository},
		PauseScheduledExpenseAction:  &scheduledexpense.PauseScheduledExpense{Repository: scheduledExpenseRepository, Today: date.Today},
		ResumeScheduledExpenseAction: &scheduledexpense.ResumeScheduledExpense{Repository: scheduledExpenseRepository, Today: date.Today},
	}
}

func NewAttachmentRepository() attachment.AttachmentRepository {
	cfg, err := awsclient.LoadDefaultConfig(
		context.Background(),
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
		GetAttachmentAction:    attachment.NewGetAttachment(repository),
		DeleteAttachmentAction: attachment.NewDeleteAttachment(repository),
	}
}

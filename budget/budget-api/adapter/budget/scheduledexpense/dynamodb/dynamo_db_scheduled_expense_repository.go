package dynamodb

import (
	"context"
	"errors"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/scheduledexpense"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

// DynamoDbScheduledExpenseRepository stores each ScheduledExpense under PK
// user_name / SK id (see budget/budget-api/docs/adr/0005-scheduled-expense-recurrence-and-generation-engine.md).
// Unlike expense/revenue's derived composite keys, user_name is used verbatim
// as the partition key, so FindAll's ownership scoping is structural — no
// separate filter expression is needed the way revenue/expense need one.
type DynamoDbScheduledExpenseRepository struct {
	TableName                  string
	Client                     *dynamodb.Client
	ScheduledExpenseIdProvider scheduledexpense.ScheduledExpenseIdProvider
	SearchTagRepository        tags.SearchTagRepository
	logger                     *logging.Logger
}

func NewDynamoDbScheduledExpenseRepository(
	tableName string,
	client *dynamodb.Client,
	idProvider scheduledexpense.ScheduledExpenseIdProvider,
	searchTagRepository tags.SearchTagRepository,
) scheduledexpense.ScheduledExpenseRepository {
	return &DynamoDbScheduledExpenseRepository{
		TableName:                  tableName,
		Client:                     client,
		ScheduledExpenseIdProvider: idProvider,
		SearchTagRepository:        searchTagRepository,
		logger:                     logging.GetLoggerInstanceForComponentByType(&DynamoDbScheduledExpenseRepository{}),
	}
}

func (repository *DynamoDbScheduledExpenseRepository) Save(ctx context.Context, se *scheduledexpense.ScheduledExpense) error {
	if se.Id == "" {
		se.Id = repository.ScheduledExpenseIdProvider.GenerateIdFor(se)
	}

	// UserName is the partition key — it must come from the authenticated
	// caller, never trusted from the struct (see revenue's Save for the same
	// enforcement, and budget-api/CLAUDE.md "Ownership enforced in two
	// layers"). CreateScheduledExpense.Execute already sets it correctly
	// before calling Save; this re-derivation is the backstop for callers
	// that don't (e.g. an Update action built from a representation that
	// never carries UserName on the wire).
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return err
	}
	se.UserName = *user.UserName

	tagKeys := make([]string, 0, len(se.Tags))
	for _, tag := range se.Tags {
		tagKeys = append(tagKeys, tag.Key)
	}

	item := map[string]types.AttributeValue{
		"user_name":   &types.AttributeValueMemberS{Value: se.UserName},
		"id":          &types.AttributeValueMemberS{Value: se.Id},
		"description": &types.AttributeValueMemberS{Value: se.Description},
		"amount":      &types.AttributeValueMemberS{Value: se.Amount.StringifyAmount()},
		"notes":       &types.AttributeValueMemberS{Value: se.Notes},
		"tag":         &types.AttributeValueMemberS{Value: strings.Join(tagKeys, ",")},
		"day":         &types.AttributeValueMemberN{Value: strconv.Itoa(se.Day)},
		"status":      &types.AttributeValueMemberS{Value: string(se.Status)},
	}
	if se.Month != nil {
		item["month"] = &types.AttributeValueMemberN{Value: strconv.Itoa(*se.Month)}
	}
	if se.EndDate != nil {
		item["end_date"] = &types.AttributeValueMemberS{Value: se.EndDate.GetIsoFormattedDate()}
	}
	if se.LastEvaluatedDate != nil {
		item["last_evaluated_date"] = &types.AttributeValueMemberS{Value: se.LastEvaluatedDate.GetIsoFormattedDate()}
	}

	_, err = repository.Client.PutItem(ctx, &dynamodb.PutItemInput{
		TableName: aws.String(repository.TableName),
		Item:      item,
	})
	return err
}

func (repository *DynamoDbScheduledExpenseRepository) FindAll(ctx context.Context) ([]scheduledexpense.ScheduledExpense, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName:              aws.String(repository.TableName),
		KeyConditionExpression: aws.String("user_name = :user_name"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_name": &types.AttributeValueMemberS{Value: *user.UserName},
		},
	}

	result, err := repository.Client.Query(ctx, input)
	if err != nil {
		repository.logger.LogErrorfFor("Error querying DynamoDB: %v", err)
		return nil, err
	}

	scheduledExpenses := make([]scheduledexpense.ScheduledExpense, 0, len(result.Items))
	for _, item := range result.Items {
		se, err := repository.fromDynamo(ctx, item)
		if err != nil {
			repository.logger.LogErrorfFor("Error processing item in FindAll: %v", err)
			continue
		}
		scheduledExpenses = append(scheduledExpenses, *se)
	}
	return scheduledExpenses, nil
}

func (repository *DynamoDbScheduledExpenseRepository) fromDynamo(ctx context.Context, item map[string]types.AttributeValue) (*scheduledexpense.ScheduledExpense, error) {
	amount, err := money.MoneyFor(item["amount"].(*types.AttributeValueMemberS).Value)
	if err != nil {
		repository.logger.LogErrorfFor("invalid data format in ScheduledExpense: %v", err)
		return nil, errors.New("invalid data format in ScheduledExpense")
	}

	dayAttr, ok := item["day"].(*types.AttributeValueMemberN)
	if !ok {
		return nil, errors.New("invalid data format in ScheduledExpense")
	}
	day, err := strconv.Atoi(dayAttr.Value)
	if err != nil {
		return nil, errors.New("invalid data format in ScheduledExpense")
	}

	searchTags, err := repository.resolveTags(ctx, item)
	if err != nil {
		return nil, err
	}

	se := &scheduledexpense.ScheduledExpense{
		Id:          item["id"].(*types.AttributeValueMemberS).Value,
		UserName:    item["user_name"].(*types.AttributeValueMemberS).Value,
		Description: item["description"].(*types.AttributeValueMemberS).Value,
		Amount:      amount,
		Notes:       item["notes"].(*types.AttributeValueMemberS).Value,
		Tags:        searchTags,
		Day:         day,
		Status:      scheduledexpense.Status(item["status"].(*types.AttributeValueMemberS).Value),
	}

	if monthAttr, ok := item["month"].(*types.AttributeValueMemberN); ok {
		if month, err := strconv.Atoi(monthAttr.Value); err == nil {
			se.Month = &month
		}
	}
	if endDateAttr, ok := item["end_date"].(*types.AttributeValueMemberS); ok {
		if d, err := date.IsoDateFor(endDateAttr.Value); err == nil {
			se.EndDate = d
		}
	}
	if lastEvalAttr, ok := item["last_evaluated_date"].(*types.AttributeValueMemberS); ok {
		if d, err := date.IsoDateFor(lastEvalAttr.Value); err == nil {
			se.LastEvaluatedDate = d
		}
	}

	return se, nil
}

// resolveTags turns the stored comma-joined tag keys into SearchTags with
// their current values, fetched from tag-api — the same key-only, resolve-on-
// read convention BudgetExpense and Revenue use. An empty "tag" attribute
// (an untagged template) resolves to the UNKNOWN sentinel.
func (repository *DynamoDbScheduledExpenseRepository) resolveTags(ctx context.Context, item map[string]types.AttributeValue) ([]tags.SearchTag, error) {
	tagAttr, ok := item["tag"].(*types.AttributeValueMemberS)
	if !ok || tagAttr.Value == "" {
		return []tags.SearchTag{tags.UnknownSentinel()}, nil
	}

	tagKeys := strings.Split(tagAttr.Value, ",")
	searchTags := make([]tags.SearchTag, 0, len(tagKeys))
	for _, tagKey := range tagKeys {
		searchTag, err := repository.SearchTagRepository.GetTagBy(ctx, tagKey)
		if err != nil {
			repository.logger.LogErrorfFor("Error getting tag: %v", err)
			return nil, errors.New("invalid tag in ScheduledExpense")
		}
		searchTags = append(searchTags, *searchTag)
	}
	return searchTags, nil
}

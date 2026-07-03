package dynamodb

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/expense"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/tags"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type DynamoDbBudgetExpenseRepository struct {
	TableName               string
	Client                  *dynamodb.Client
	BudgetExpenseIdProvider expense.BudgetExpenseIdProvider
	SearchTagRepository     tags.SearchTagRepository
	logger                  *logging.Logger
	EventBus                *expense.InternalEventBus
}

func NewDynamoDbBudgetExpenseRepository(
	TableName string,
	Client *dynamodb.Client,
	BudgetExpenseIdProvider expense.BudgetExpenseIdProvider,
	SearchTagRepository tags.SearchTagRepository,
	EventBus *expense.InternalEventBus) expense.BudgetExpenseRepository {
	return &DynamoDbBudgetExpenseRepository{
		TableName:               TableName,
		Client:                  Client,
		BudgetExpenseIdProvider: BudgetExpenseIdProvider,
		SearchTagRepository:     SearchTagRepository,
		logger:                  logging.GetLoggerInstanceForComponentByType(&DynamoDbBudgetExpenseRepository{}),
		EventBus:                EventBus,
	}
}

func (repository *DynamoDbBudgetExpenseRepository) FindFor(ctx context.Context, budgetExpenseId expense.BudgetExpenseId) (*expense.BudgetExpense, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error to get a valid user in the context: %v", err)
		return nil, err
	}

	pk, range_key, err := dynamoDbKeysFrom(budgetExpenseId)
	if err != nil {
		repository.logger.LogErrorfFor("Error dynamodb key generation: %v", err)
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName: aws.String(repository.TableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":pk":        &types.AttributeValueMemberS{Value: pk},
			":range_key": &types.AttributeValueMemberS{Value: range_key},
			":user_name": &types.AttributeValueMemberS{Value: *user.UserName},
		},
		KeyConditionExpression: aws.String("pk =:pk AND range_key =:range_key"),
		FilterExpression:       aws.String("user_name = :user_name"),
	}
	result, err := repository.Client.Query(ctx, input)
	if err != nil {
		repository.logger.LogErrorfFor("Error querying DynamoDB: %v", err)
		return nil, err
	}

	if len(result.Items) == 0 {
		repository.logger.LogErrorfFor("No budget expense found for id: %s", budgetExpenseId)
		return nil, errors.New("BudgetExpense not found")
	}

	item := result.Items[0]

	// FindFor is the ownership pre-check for update/delete, never a standalone
	// read, so it deliberately drops the deleted-tag flag and never reclassifies:
	// an async re-Save of the pre-mutation snapshot would resurrect a just-deleted
	// expense or clobber an in-flight update. Only FindByDateRange reclassifies.
	budgetExpense, _, err := repository.fromDynamo(ctx, item)
	return budgetExpense, err
}

func (repository *DynamoDbBudgetExpenseRepository) FindByDateRange(ctx context.Context, start date.Date, end date.Date, searchTags []tags.SearchTagKey) ([]expense.BudgetExpense, error) {
	budgetExpenses := make([]expense.BudgetExpense, 0)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return nil, err
	}
	idProvider, _ := repository.BudgetExpenseIdProvider.(*DynamoDbBudgetExpenseIdProvider)

	filterExpression := "user_name = :user_name AND transaction_date BETWEEN :start_date AND :end_date"
	expressionAttributeValues := map[string]types.AttributeValue{
		":user_name":  &types.AttributeValueMemberS{Value: *user.UserName},
		":start_date": &types.AttributeValueMemberS{Value: start.GetIsoFormattedDate()},
		":end_date":   &types.AttributeValueMemberS{Value: end.GetIsoFormattedDate()},
	}
	tagConditions := make([]string, 0, len(searchTags))
	if len(searchTags) > 0 {
		for index, searchTag := range searchTags {
			placeholder := fmt.Sprintf(":tag_%d", index)
			tagConditions = append(tagConditions, "contains(tag, "+placeholder+")")
			expressionAttributeValues[placeholder] = &types.AttributeValueMemberS{Value: searchTag}
		}
		filterExpression += " AND (" + strings.Join(tagConditions, " OR ") + ")"
	}

	partitionKeys := repository.partitionKeysForDateRangeAndUser(idProvider, start, end, *user.UserName)
	for _, pk := range partitionKeys {
		expressionAttributeValues[":pk"] = &types.AttributeValueMemberS{Value: pk}
		input := &dynamodb.QueryInput{
			TableName:                 aws.String(repository.TableName),
			ExpressionAttributeValues: expressionAttributeValues,
			FilterExpression:          aws.String(filterExpression),
			KeyConditionExpression:    aws.String("pk = :pk"),
		}

		items, err := repository.Client.Query(ctx, input)
		if err != nil {
			return nil, err
		}
		for _, item := range items.Items {
			budgetExpense, needsReclassification, err := repository.fromDynamo(ctx, item)
			if err != nil {
				repository.logger.LogErrorfFor("Error processing item in FindByDateRange: %v", err)
			} else {
				budgetExpenses = append(budgetExpenses, *budgetExpense)
				if needsReclassification {
					repository.publishReclassification(ctx, budgetExpense)
				}
			}
		}
	}

	return budgetExpenses, nil
}

// fromDynamo decodes a DynamoDB item into a BudgetExpense and reports whether it
// carries a deleted-tag reference (a stored key that resolved to the UNKNOWN
// sentinel). It has no side effect: the caller decides whether to durably
// reclassify, so the write is confined to genuine read paths (FindByDateRange)
// and never to FindFor's update/delete pre-check.
func (repository *DynamoDbBudgetExpenseRepository) fromDynamo(ctx context.Context, item map[string]types.AttributeValue) (*expense.BudgetExpense, bool, error) {
	date, err := date.IsoDateFor(item["transaction_date"].(*types.AttributeValueMemberS).Value)
	if err != nil {
		repository.logger.LogErrorfFor("Error parsing transaction_date: %v", err)
		return nil, false, errors.New("invalid data format in BudgetExpense")
	}

	moneyAmount, err := money.MoneyFor(item["amount"].(*types.AttributeValueMemberS).Value)
	if err != nil {
		repository.logger.LogErrorfFor("invalid data format in BudgetExpense: %v", err)
		return nil, false, errors.New("invalid data format in BudgetExpense")
	}

	tagKeys := strings.Split(item["tag"].(*types.AttributeValueMemberS).Value, ",")
	searchTags := make([]tags.SearchTag, 0, len(tagKeys))
	unknownKey := tags.UnknownSentinel().Key
	hasDeletedTagReference := false

	for _, tagKey := range tagKeys {
		searchTag, err := repository.SearchTagRepository.GetTagBy(ctx, tagKey)
		if err != nil {
			repository.logger.LogErrorfFor("Error getting tag: %v", err)
			return nil, false, errors.New("invalid tag in BudgetExpense")
		}
		// A stored key that resolves to the UNKNOWN sentinel while the stored key
		// itself was not UNKNOWN is a reference to a tag deleted in tag-api, and
		// must be durably reclassified. A record stored with the UNKNOWN sentinel
		// as its key (the default applied at create time) resolves to UNKNOWN too
		// but is not a deletion, so it is excluded by the second comparison.
		if searchTag.Key == unknownKey && tagKey != unknownKey {
			hasDeletedTagReference = true
		}
		searchTags = append(searchTags, *searchTag)
	}

	budgetExpense := &expense.BudgetExpense{
		Id:       item["budget_id"].(*types.AttributeValueMemberS).Value,
		UserName: item["user_name"].(*types.AttributeValueMemberS).Value,
		Date:     *date,
		Amount:   moneyAmount,
		Note:     item["note"].(*types.AttributeValueMemberS).Value,
		Tags:     searchTags,
	}
	return budgetExpense, hasDeletedTagReference, nil
}

// publishReclassification hands a copy of the just-read expense to the
// reclassification bus so the service layer can durably rewrite its deleted-tag
// references to UNKNOWN. It sends a clone — never the pointer returned to the
// caller — so the async listener and the read caller never share mutable state.
// The clone's tags are de-duplicated so a record whose every tag was deleted
// persists a single UNKNOWN rather than repeated sentinels.
func (repository *DynamoDbBudgetExpenseRepository) publishReclassification(ctx context.Context, source *expense.BudgetExpense) {
	clone := *source
	clone.Tags = distinctTags(source.Tags)
	repository.EventBus.Publish(expense.InternalEvent{Payload: &clone, Ctx: ctx})
}

func distinctTags(searchTags []tags.SearchTag) []tags.SearchTag {
	seen := make(map[string]bool, len(searchTags))
	distinct := make([]tags.SearchTag, 0, len(searchTags))
	for _, searchTag := range searchTags {
		if seen[searchTag.Key] {
			continue
		}
		seen[searchTag.Key] = true
		distinct = append(distinct, searchTag)
	}
	return distinct
}

func (repository *DynamoDbBudgetExpenseRepository) partitionKeysForDateRangeAndUser(idProvider *DynamoDbBudgetExpenseIdProvider, start, end date.Date, userName security.UserName) []string {
	// Generate partition keys for each month in the date range
	var partitionKeys []string
	currentMonth := start.GetMonth()
	currentYear := start.GetYear()
	lastMonth := end.GetMonth()
	lastYear := end.GetYear()
	for currentMonth < lastMonth || currentYear <= lastYear {
		partitionKeys = append(partitionKeys, idProvider.partitionKeyFrom(currentYear, currentMonth, userName))
		currentMonth++
		if currentMonth > 12 {
			currentMonth = 1
			currentYear++
		}
	}
	return partitionKeys
}

func (repository *DynamoDbBudgetExpenseRepository) Save(ctx context.Context, budgetExpense *expense.BudgetExpense) error {
	isNew := false
	if budgetExpense.Id == "" {
		isNew = true
		budgetExpense.Id = repository.BudgetExpenseIdProvider.GenerateIdFor(budgetExpense)
	}
	pk, range_key, err := dynamoDbKeysFrom(budgetExpense.Id)
	if err != nil {
		return err
	}

	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	budgetExpense.UserName = *user.UserName

	tagKeys := make([]string, 0, len(budgetExpense.Tags))
	for _, tag := range budgetExpense.Tags {
		tagKeys = append(tagKeys, tag.Key)
	}

	queryInput := &dynamodb.PutItemInput{
		TableName: &repository.TableName,

		Item: map[string]types.AttributeValue{
			"pk":               &types.AttributeValueMemberS{Value: pk},
			"range_key":        &types.AttributeValueMemberS{Value: range_key},
			"user_name":        &types.AttributeValueMemberS{Value: budgetExpense.UserName},
			"budget_id":        &types.AttributeValueMemberS{Value: budgetExpense.Id},
			"transaction_date": &types.AttributeValueMemberS{Value: budgetExpense.Date.GetIsoFormattedDate()},
			"amount":           &types.AttributeValueMemberS{Value: budgetExpense.Amount.StringifyAmount()},
			"note":             &types.AttributeValueMemberS{Value: budgetExpense.Note},
			"tag":              &types.AttributeValueMemberS{Value: strings.Join(tagKeys, ",")},
		},
	}

	if !isNew {
		// to update, we need to make sure that the user name matches
		queryInput.ConditionExpression = aws.String("user_name = :user_name")
		queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
			":user_name": &types.AttributeValueMemberS{Value: budgetExpense.UserName},
		}
	}

	_, err = repository.Client.PutItem(ctx, queryInput)

	// Implementation to save a budget expense to DynamoDB
	return err
}

func (repository *DynamoDbBudgetExpenseRepository) Delete(ctx context.Context, idBudgetExpense expense.BudgetExpenseId) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	// Implementation to delete a budget expense from DynamoDB
	pk, range_key, err := dynamoDbKeysFrom(idBudgetExpense)
	if err != nil {
		return err
	}

	_, err = repository.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: &repository.TableName,
		Key: map[string]types.AttributeValue{
			"pk":        &types.AttributeValueMemberS{Value: pk},
			"range_key": &types.AttributeValueMemberS{Value: range_key},
		},
		ConditionExpression: aws.String("user_name = :user_name"),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_name": &types.AttributeValueMemberS{Value: *user.UserName},
		},
	})

	if err != nil {
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			return errors.New("user name does not match")
		} else {
			return errors.New("unexpected error deleting budget expense")
		}
	}
	return err
}

func dynamoDbKeysFrom(budgetExpenseId expense.BudgetExpenseId) (string, string, error) {
	kyes := strings.Split(string(budgetExpenseId), "-")
	if len(kyes) != 2 {
		return "", "", errors.New("invalid budget expense id format: " + string(budgetExpenseId))
	}
	pk := kyes[0]
	range_key := kyes[1]

	return pk, range_key, nil
}

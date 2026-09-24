package dynamodb

import (
	"context"
	"errors"
	"fmt"
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

	// scanPageLimit caps each FindAllActive Scan page; 0 means DynamoDB's
	// default (1MB pages). Tests lower it to exercise pagination.
	scanPageLimit int32
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
	tagNames := make(map[string]types.AttributeValue, len(se.Tags))
	for _, tag := range se.Tags {
		tagKeys = append(tagKeys, tag.Key)
		tagNames[tag.Key] = &types.AttributeValueMemberS{Value: tag.Value}
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
	// tag_names keeps each tag's display name as of this save, for the
	// generation engine: it runs with no user token, so it can't resolve names
	// through tag-api the way user-facing reads do (#55). Reads for the UI
	// still resolve live names and ignore this attribute.
	if len(tagNames) > 0 {
		item["tag_names"] = &types.AttributeValueMemberM{Value: tagNames}
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

// FindFor is a point lookup on the full primary key (PK=user_name from ctx,
// SK=id), so it uses GetItem rather than Query — unlike FindAll, which
// necessarily queries (PK only). Returns (nil, nil) on a miss: no item at
// that key under this user's own partition, whether because the id doesn't
// exist at all or because it belongs to a different user — the two are
// indistinguishable by construction and neither is an error condition (see
// the ScheduledExpenseRepository port doc).
func (repository *DynamoDbScheduledExpenseRepository) FindFor(ctx context.Context, id scheduledexpense.ScheduledExpenseId) (*scheduledexpense.ScheduledExpense, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return nil, err
	}

	result, err := repository.Client.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(repository.TableName),
		Key: map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{Value: *user.UserName},
			"id":        &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		repository.logger.LogErrorfFor("Error getting item from DynamoDB: %v", err)
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	return repository.fromDynamo(ctx, repository.resolveTags, result.Item)
}

// Delete removes the item at (PK=user_name from ctx, SK=id). Since the owner
// is part of the key, another user's row is unreachable by construction; the
// attribute_exists(id) condition makes a missing row a
// ErrScheduledExpenseNotFound rather than DynamoDB's silent no-op success.
func (repository *DynamoDbScheduledExpenseRepository) Delete(ctx context.Context, id scheduledexpense.ScheduledExpenseId) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return err
	}

	_, err = repository.Client.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(repository.TableName),
		Key: map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{Value: *user.UserName},
			"id":        &types.AttributeValueMemberS{Value: id},
		},
		ConditionExpression: aws.String("attribute_exists(id)"),
	})
	if err != nil {
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			return scheduledexpense.ErrScheduledExpenseNotFound
		}
		repository.logger.LogErrorfFor("Error deleting item from DynamoDB: %v", err)
		return err
	}
	return nil
}

// UpdateStatus is a targeted UpdateItem on status + last_evaluated_date only —
// not a Save of a FindFor-resolved record, which would rewrite every attribute
// (tag keys included) and clobber a concurrent edit. attribute_exists(id)
// both maps a missing row to ErrScheduledExpenseNotFound and stops UpdateItem
// from upserting a stub row. "status" is a DynamoDB reserved word, hence the
// #status alias.
func (repository *DynamoDbScheduledExpenseRepository) UpdateStatus(ctx context.Context, id scheduledexpense.ScheduledExpenseId, status scheduledexpense.Status, lastEvaluatedDate date.Date) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return err
	}

	_, err = repository.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(repository.TableName),
		Key: map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{Value: *user.UserName},
			"id":        &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:    aws.String("SET #status = :status, last_evaluated_date = :last_evaluated_date"),
		ConditionExpression: aws.String("attribute_exists(id)"),
		ExpressionAttributeNames: map[string]string{
			"#status": "status",
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":status":              &types.AttributeValueMemberS{Value: string(status)},
			":last_evaluated_date": &types.AttributeValueMemberS{Value: lastEvaluatedDate.GetIsoFormattedDate()},
		},
	})
	if err != nil {
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			return scheduledexpense.ErrScheduledExpenseNotFound
		}
		repository.logger.LogErrorfFor("Error updating scheduled expense status in DynamoDB: %v", err)
		return err
	}
	return nil
}

// FindAllActive scans the whole table — every user's partition — for ACTIVE
// definitions, following LastEvaluatedKey across pages. Tags come from
// storedTags, never tag-api. A row that can't be decoded is logged and
// skipped so one bad definition can't stop generation for everyone.
func (repository *DynamoDbScheduledExpenseRepository) FindAllActive(ctx context.Context) ([]scheduledexpense.ScheduledExpense, error) {
	input := &dynamodb.ScanInput{
		TableName:                aws.String(repository.TableName),
		FilterExpression:         aws.String("#status = :active"),
		ExpressionAttributeNames: map[string]string{"#status": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":active": &types.AttributeValueMemberS{Value: string(scheduledexpense.StatusActive)},
		},
	}
	if repository.scanPageLimit > 0 {
		input.Limit = aws.Int32(repository.scanPageLimit)
	}

	result := []scheduledexpense.ScheduledExpense{}
	paginator := dynamodb.NewScanPaginator(repository.Client, input)
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			repository.logger.LogErrorfFor("Error scanning active scheduled expenses: %v", err)
			return nil, err
		}
		for _, item := range page.Items {
			se, err := repository.decodeForGeneration(ctx, item)
			if err != nil {
				repository.logger.LogErrorfFor("Skipping undecodable scheduled expense row: %v", err)
				continue
			}
			result = append(result, *se)
		}
	}
	return result, nil
}

// decodeForGeneration decodes one scanned row, turning a panic from
// fromDynamo's unchecked attribute assertions (a row missing an attribute)
// into an error, so the scan skips that row instead of losing the whole run.
func (repository *DynamoDbScheduledExpenseRepository) decodeForGeneration(ctx context.Context, item map[string]types.AttributeValue) (se *scheduledexpense.ScheduledExpense, err error) {
	defer func() {
		if r := recover(); r != nil {
			se, err = nil, fmt.Errorf("malformed row: %v", r)
		}
	}()
	return repository.fromDynamo(ctx, storedTags, item)
}

// AdvanceLastEvaluatedDate sets only last_evaluated_date on (PK=user_name
// from ctx, SK=id). The condition requires the row to still exist and be
// ACTIVE: a definition deleted mid-run isn't upserted back as a stub row, and
// one paused mid-run keeps the date its pause stamped.
func (repository *DynamoDbScheduledExpenseRepository) AdvanceLastEvaluatedDate(ctx context.Context, id scheduledexpense.ScheduledExpenseId, lastEvaluatedDate date.Date) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return err
	}

	_, err = repository.Client.UpdateItem(ctx, &dynamodb.UpdateItemInput{
		TableName: aws.String(repository.TableName),
		Key: map[string]types.AttributeValue{
			"user_name": &types.AttributeValueMemberS{Value: *user.UserName},
			"id":        &types.AttributeValueMemberS{Value: id},
		},
		UpdateExpression:         aws.String("SET last_evaluated_date = :last_evaluated_date"),
		ConditionExpression:      aws.String("attribute_exists(id) AND #status = :active"),
		ExpressionAttributeNames: map[string]string{"#status": "status"},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":last_evaluated_date": &types.AttributeValueMemberS{Value: lastEvaluatedDate.GetIsoFormattedDate()},
			":active":              &types.AttributeValueMemberS{Value: string(scheduledexpense.StatusActive)},
		},
	})
	if err != nil {
		var conditionalCheckFailedException *types.ConditionalCheckFailedException
		if errors.As(err, &conditionalCheckFailedException) {
			return scheduledexpense.ErrScheduledExpenseNotFound
		}
		repository.logger.LogErrorfFor("Error advancing scheduled expense last evaluated date in DynamoDB: %v", err)
		return err
	}
	return nil
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
		se, err := repository.fromDynamo(ctx, repository.resolveTags, item)
		if err != nil {
			repository.logger.LogErrorfFor("Error processing item in FindAll: %v", err)
			continue
		}
		scheduledExpenses = append(scheduledExpenses, *se)
	}
	return scheduledExpenses, nil
}

// tagReader turns a stored item's tags into SearchTags: resolveTags (live
// names from tag-api, for user-facing reads) or storedTags (names saved on the
// row, for the token-less generation engine).
type tagReader func(ctx context.Context, item map[string]types.AttributeValue) ([]tags.SearchTag, error)

func (repository *DynamoDbScheduledExpenseRepository) fromDynamo(ctx context.Context, readTags tagReader, item map[string]types.AttributeValue) (*scheduledexpense.ScheduledExpense, error) {
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

	searchTags, err := readTags(ctx, item)
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

// storedTags reads tag keys and the names stored at save time, without calling
// tag-api. Rows saved before names were stored yield empty names.
func storedTags(_ context.Context, item map[string]types.AttributeValue) ([]tags.SearchTag, error) {
	tagAttr, ok := item["tag"].(*types.AttributeValueMemberS)
	if !ok || tagAttr.Value == "" {
		return []tags.SearchTag{tags.UnknownSentinel()}, nil
	}

	names := map[string]types.AttributeValue{}
	if namesAttr, ok := item["tag_names"].(*types.AttributeValueMemberM); ok {
		names = namesAttr.Value
	}

	tagKeys := strings.Split(tagAttr.Value, ",")
	searchTags := make([]tags.SearchTag, 0, len(tagKeys))
	for _, tagKey := range tagKeys {
		searchTag := tags.SearchTag{Key: tagKey}
		if name, ok := names[tagKey].(*types.AttributeValueMemberS); ok {
			searchTag.Value = name.Value
		}
		searchTags = append(searchTags, searchTag)
	}
	return searchTags, nil
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

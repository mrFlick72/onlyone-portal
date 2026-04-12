package dynamodb

import (
	"context"
	"errors"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrflick72/budget/budget-api/domain/budget/revenue"
	"github.com/mrflick72/budget/budget-api/domain/money"
	"github.com/mrflick72/budget/budget-api/domain/time/date"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type DynamoDbRevenueRepository struct {
	TableName         string
	Client            *dynamodb.Client
	RevenueIdProvider revenue.RevenueIdProvider
	logger            *logging.Logger
}

func NewDynamoDbRevenueRepository(
	TableName string,
	Client *dynamodb.Client,
	RevenueIdProvider revenue.RevenueIdProvider) revenue.RevenueRepository {
	return &DynamoDbRevenueRepository{
		TableName:         TableName,
		Client:            Client,
		RevenueIdProvider: RevenueIdProvider,
		logger:            logging.GetLoggerInstanceForComponentByType(&DynamoDbRevenueRepository{}),
	}
}

func (repository *DynamoDbRevenueRepository) FindFor(ctx context.Context, revenueId revenue.RevenueId) (*revenue.Revenue, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error to get a valid user in the context: %v", err)
		return nil, err
	}

	pk, range_key, err := dynamoDbKeysFrom(revenueId)
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
	result, err := repository.Client.Query(context.TODO(), input)
	if err != nil {
		repository.logger.LogErrorfFor("Error querying DynamoDB: %v", err)
		return nil, err
	}

	if len(result.Items) == 0 {
		repository.logger.LogErrorfFor("No revenue found for id: %s", revenueId)
		return nil, errors.New("Revenue not found")
	}

	return repository.fromDynamo(result.Items[0])
}

func (repository *DynamoDbRevenueRepository) fromDynamo(item map[string]types.AttributeValue) (*revenue.Revenue, error) {
	d, err := date.IsoDateFor(item["transaction_date"].(*types.AttributeValueMemberS).Value)
	if err != nil {
		repository.logger.LogErrorfFor("Error parsing transaction_date: %v", err)
		return nil, errors.New("invalid data format in Revenue")
	}

	amount, err := money.MoneyFor(item["amount"].(*types.AttributeValueMemberS).Value)
	if err != nil {
		repository.logger.LogErrorfFor("invalid data format in Revenue: %v", err)
		return nil, errors.New("invalid data format in Revenue")
	}

	return &revenue.Revenue{
		Id:       revenue.RevenueId(item["budget_id"].(*types.AttributeValueMemberS).Value),
		UserName: item["user_name"].(*types.AttributeValueMemberS).Value,
		Date:     *d,
		Amount:   amount,
		Note:     item["note"].(*types.AttributeValueMemberS).Value,
	}, nil
}

func (repository *DynamoDbRevenueRepository) FindByDateRange(ctx context.Context, start date.Date, end date.Date) ([]revenue.Revenue, error) {
	revenues := make([]revenue.Revenue, 0)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		repository.logger.LogErrorfFor("Error getting current user: %v", err)
		return nil, err
	}

	idProvider, _ := repository.RevenueIdProvider.(*DynamoDbRevenueIdProvider)
	partitionKeys := repository.partitionKeysForDateRangeAndUser(idProvider, start, end, *user.UserName)

	for _, pk := range partitionKeys {
		input := &dynamodb.QueryInput{
			TableName: aws.String(repository.TableName),
			ExpressionAttributeValues: map[string]types.AttributeValue{
				":user_name":  &types.AttributeValueMemberS{Value: *user.UserName},
				":start_date": &types.AttributeValueMemberS{Value: start.GetIsoFormattedDate()},
				":end_date":   &types.AttributeValueMemberS{Value: end.GetIsoFormattedDate()},
				":pk":         &types.AttributeValueMemberS{Value: pk},
			},
			FilterExpression:       aws.String("user_name = :user_name AND transaction_date BETWEEN :start_date AND :end_date"),
			KeyConditionExpression: aws.String("pk = :pk"),
		}

		items, err := repository.Client.Query(context.TODO(), input)
		if err != nil {
			return nil, err
		}
		for _, item := range items.Items {
			r, err := repository.fromDynamo(item)
			if err != nil {
				repository.logger.LogErrorfFor("Error processing item in FindByDateRange: %v", err)
			} else {
				revenues = append(revenues, *r)
			}
		}
	}

	return revenues, nil
}

func (repository *DynamoDbRevenueRepository) partitionKeysForDateRangeAndUser(idProvider *DynamoDbRevenueIdProvider, start, end date.Date, userName string) []string {
	var partitionKeys []string
	for year := start.GetYear(); year <= end.GetYear(); year++ {
		partitionKeys = append(partitionKeys, idProvider.partitionKeyFrom(year, userName))
	}
	return partitionKeys
}

func (repository *DynamoDbRevenueRepository) Save(ctx context.Context, revenueToSave *revenue.Revenue) error {
	isNew := false
	if revenueToSave.Id == "" {
		isNew = true
		revenueToSave.Id = repository.RevenueIdProvider.GenerateIdFor(revenueToSave)
	}
	pk, range_key, err := dynamoDbKeysFrom(revenueToSave.Id)
	if err != nil {
		return err
	}

	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}
	revenueToSave.UserName = *user.UserName

	queryInput := &dynamodb.PutItemInput{
		TableName: &repository.TableName,
		Item: map[string]types.AttributeValue{
			"pk":               &types.AttributeValueMemberS{Value: pk},
			"range_key":        &types.AttributeValueMemberS{Value: range_key},
			"user_name":        &types.AttributeValueMemberS{Value: revenueToSave.UserName},
			"budget_id":        &types.AttributeValueMemberS{Value: string(revenueToSave.Id)},
			"transaction_date": &types.AttributeValueMemberS{Value: revenueToSave.Date.GetIsoFormattedDate()},
			"amount":           &types.AttributeValueMemberS{Value: revenueToSave.Amount.StringifyAmount()},
			"note":             &types.AttributeValueMemberS{Value: revenueToSave.Note},
		},
	}

	if !isNew {
		queryInput.ConditionExpression = aws.String("user_name = :user_name")
		queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
			":user_name": &types.AttributeValueMemberS{Value: revenueToSave.UserName},
		}
	}

	_, err = repository.Client.PutItem(ctx, queryInput)
	return err
}

func (repository *DynamoDbRevenueRepository) Delete(ctx context.Context, revenueId revenue.RevenueId) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	pk, range_key, err := dynamoDbKeysFrom(revenueId)
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
		}
		return errors.New("unexpected error deleting revenue")
	}
	return nil
}

func dynamoDbKeysFrom(revenueId revenue.RevenueId) (string, string, error) {
	keys := strings.Split(string(revenueId), "-")
	if len(keys) != 2 {
		return "", "", errors.New("invalid revenue id format: " + string(revenueId))
	}
	return keys[0], keys[1], nil
}

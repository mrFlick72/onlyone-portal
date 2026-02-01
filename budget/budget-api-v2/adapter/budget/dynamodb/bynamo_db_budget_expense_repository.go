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
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

type DynamoDbBudgetExpenseRepository struct {
	TableName               string
	Client                  *dynamodb.Client
	BudgetExpenseIdProvider expense.BudgetExpenseIdProvider
}

func (repository *DynamoDbBudgetExpenseRepository) FindByDateRange(ctx *context.Context, start date.Date, end date.Date, searchTags []tags.SearchTagKey) (*[]expense.BudgetExpense, error) {
	// Implementation to interact with DynamoDB and retrieve budget expenses by date range and search tags
	budgetExpenses := make([]expense.BudgetExpense, 0)
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		fmt.Println("Error getting current user:", err)
		return nil, err
	}
	// For simplicity, returning nil. Actual implementation would query DynamoDB.
	input := &dynamodb.QueryInput{
		TableName: aws.String(repository.TableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_name":  &types.AttributeValueMemberS{Value: *user.UserName},
			":start_date": &types.AttributeValueMemberS{Value: start.GetIsoFormattedDate()},
			":end_date":   &types.AttributeValueMemberS{Value: end.GetIsoFormattedDate()},
		},
		FilterExpression:       aws.String("user_name = :user_name AND transaction_date BETWEEN :start_date AND :end_date"),
		KeyConditionExpression: aws.String("pk = :pk"),
	}

	items, err := repository.Client.Query(context.TODO(), input)
	if err != nil {
		return nil, err
	}
	for _, item := range items.Items {

		date, err := date.IsoDateFor(item["transaction_date"].(*types.AttributeValueMemberS).Value)
		if err != nil {
			return nil, errors.New("invalid data format in BudgetExpense")
		}

		moneyAmount, err := money.MoneyFor(item["amount"].(*types.AttributeValueMemberN).Value)
		if err != nil {
			return nil, errors.New("invalid data format in BudgetExpense")
		}

		budgetExpenses = append(budgetExpenses, expense.BudgetExpense{
			Id:       expense.BudgetExpenseId(item["budget_id"].(*types.AttributeValueMemberS).Value),
			UserName: item["user_name"].(*types.AttributeValueMemberS).Value,
			Date:     *date,
			Amount:   *moneyAmount,
			Note:     item["note"].(*types.AttributeValueMemberS).Value,
			Tag:      item["tag"].(*types.AttributeValueMemberS).Value,
		})
		// Process each item and convert to BudgetExpense
	}

	// Placeholder return
	return &budgetExpenses, nil
}

func (repository *DynamoDbBudgetExpenseRepository) Save(ctx *context.Context, budgetExpense *expense.BudgetExpense) error {
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

	queryInput := &dynamodb.PutItemInput{
		TableName: &repository.TableName,

		Item: map[string]types.AttributeValue{
			"pk":               &types.AttributeValueMemberS{Value: pk},
			"range_key":        &types.AttributeValueMemberS{Value: range_key},
			"user_name":        &types.AttributeValueMemberS{Value: budgetExpense.UserName},
			"budget_id":        &types.AttributeValueMemberS{Value: string(budgetExpense.Id)},
			"transaction_date": &types.AttributeValueMemberS{Value: budgetExpense.Date.GetIsoFormattedDate()},
			"amount":           &types.AttributeValueMemberN{Value: budgetExpense.Amount.StringifyAmount()},
			"note":             &types.AttributeValueMemberS{Value: budgetExpense.Note},
			"tag":              &types.AttributeValueMemberS{Value: budgetExpense.Tag},
		},
	}

	if !isNew {
		// to update, we need to make sure that the user name matches
		queryInput.ConditionExpression = aws.String("user_name = :user_name")
		queryInput.ExpressionAttributeValues = map[string]types.AttributeValue{
			":user_name": &types.AttributeValueMemberS{Value: budgetExpense.UserName},
		}
	}

	_, err = repository.Client.PutItem(*ctx, queryInput)

	// Implementation to save a budget expense to DynamoDB
	return err
}

func (repository *DynamoDbBudgetExpenseRepository) Delete(ctx *context.Context, idBudgetExpense expense.BudgetExpenseId) error {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	// Implementation to delete a budget expense from DynamoDB
	pk, range_key, err := dynamoDbKeysFrom(idBudgetExpense)
	if err != nil {
		return err
	}

	_, err = repository.Client.DeleteItem(*ctx, &dynamodb.DeleteItemInput{
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

func (repository *DynamoDbBudgetExpenseRepository) FindFor(ctx *context.Context, budgetExpenseId expense.BudgetExpenseId) (*expense.BudgetExpense, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	pk, range_key, err := dynamoDbKeysFrom(budgetExpenseId)
	if err != nil {
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
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("BudgetExpense not found")
	}

	item := result.Items[0]

	date, err := date.IsoDateFor(item["transaction_date"].(*types.AttributeValueMemberS).Value)
	if err != nil {
		return nil, errors.New("invalid data format in BudgetExpense")
	}

	moneyAmount, err := money.MoneyFor(item["amount"].(*types.AttributeValueMemberN).Value)
	if err != nil {
		return nil, errors.New("invalid data format in BudgetExpense")
	}

	return &expense.BudgetExpense{
		Id:       expense.BudgetExpenseId(item["budget_id"].(*types.AttributeValueMemberS).Value),
		UserName: item["user_name"].(*types.AttributeValueMemberS).Value,
		Date:     *date,
		Amount:   *moneyAmount,
		Note:     item["note"].(*types.AttributeValueMemberS).Value,
		Tag:      item["tag"].(*types.AttributeValueMemberS).Value,
	}, nil
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

package dynamodb

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/aws"
	aws_config "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/awsclient"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/config"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/logging"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

var configurationManager = config.GetConfigurationManagerInstance()

type TagDynamoDBRepository struct {
	Client    *dynamodb.Client
	TableName string
}

func NewTagDynamoDBRepository() *TagDynamoDBRepository {
	var logger = logging.GetLoggerInstanceForComponentByTypeName("dynamodb.TagDynamoDBRepository")

	cfg, err := awsclient.LoadDefaultConfig(
		context.TODO(),
		aws_config.WithRegion("eu-central-1"),
	)
	dynamoDbEndpoint := configurationManager.GetConfigFor("aws.endpoint.dynamodb")

	if dynamoDbEndpoint != "" {
		logger.LogInfofFor("Using local DynamoDB endpoint: %s", dynamoDbEndpoint)
		cfg.BaseEndpoint = aws.String("http://localhost:4566")
	} else {
		logger.LogInfofFor("Using AWS DynamoDB service without local endpoint")
	}

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}
	return &TagDynamoDBRepository{
		Client:    dynamodb.NewFromConfig(cfg),
		TableName: configurationManager.GetConfigFor("tags.dynamo-db.table-name"),
	}
}

func (r *TagDynamoDBRepository) SaveTag(ctx context.Context, tag *domain.Tag) error {
	// Implementation for saving a tag to DynamoDB

	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item: map[string]types.AttributeValue{
			"search_tag_key":   &types.AttributeValueMemberS{Value: tag.Key},
			"search_tag_value": &types.AttributeValueMemberS{Value: tag.Value},
			"user_name":        &types.AttributeValueMemberS{Value: *user.UserName},
		},
	}
	_, err = r.Client.PutItem(ctx, input)

	return err
}

func (r *TagDynamoDBRepository) GetTagBy(ctx context.Context, key string) (*domain.Tag, error) {
	// Implementation for retrieving a tag by key from DynamoDB
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}

	input := &dynamodb.QueryInput{
		TableName: aws.String(r.TableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":user_name":      &types.AttributeValueMemberS{Value: *user.UserName},
			":search_tag_key": &types.AttributeValueMemberS{Value: key},
		},
		KeyConditionExpression: aws.String("user_name =:user_name AND search_tag_key =:search_tag_key"),
	}
	result, err := r.Client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	if len(result.Items) == 0 {
		return nil, errors.New("tag not found") // Tag not found
	}

	item := result.Items[0]
	return &domain.Tag{
		Key:   item["search_tag_key"].(*types.AttributeValueMemberS).Value,
		Value: item["search_tag_value"].(*types.AttributeValueMemberS).Value,
	}, nil
}

func (r *TagDynamoDBRepository) FindAllTags(ctx context.Context) ([]domain.Tag, error) {
	// Implementation for retrieving all tags from DynamoDB
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.TableName),
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: *user.UserName},
		},
		KeyConditionExpression: aws.String("user_name = :username"),
	}
	result, err := r.Client.Query(ctx, input)
	if err != nil {
		return nil, err
	}

	var tags []domain.Tag
	for _, item := range result.Items {
		tags = append(tags, domain.Tag{
			Key:   item["search_tag_key"].(*types.AttributeValueMemberS).Value,
			Value: item["search_tag_value"].(*types.AttributeValueMemberS).Value,
		})
	}

	return tags, nil
}

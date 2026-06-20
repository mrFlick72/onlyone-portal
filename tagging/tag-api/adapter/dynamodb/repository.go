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

	item := map[string]types.AttributeValue{
		"search_tag_key":   &types.AttributeValueMemberS{Value: tag.Key},
		"search_tag_value": &types.AttributeValueMemberS{Value: tag.Value},
		"user_name":        &types.AttributeValueMemberS{Value: *user.UserName},
	}

	// scope is optional: only set the attribute when present, so a tag
	// saved without one stays absent from the attribute entirely, the
	// same as every tag written before this field existed.
	if normalizedScope := domain.NormalizeScope(tag.Scope); normalizedScope != "" {
		item["scope"] = &types.AttributeValueMemberS{Value: normalizedScope}
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item:      item,
	}
	_, err = r.Client.PutItem(ctx, input)

	return err
}

// scopeFrom reads the optional "scope" attribute off a DynamoDB item,
// returning "" when the item predates the Scope field or never had one set.
func scopeFrom(item map[string]types.AttributeValue) string {
	if attr, ok := item["scope"].(*types.AttributeValueMemberS); ok {
		return attr.Value
	}
	return ""
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
		Scope: scopeFrom(item),
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
			Scope: scopeFrom(item),
		})
	}

	return tags, nil
}

// FindTagsByScope queries the same partition as FindAllTags (by user_name)
// and applies a FilterExpression on the normalized scope attribute. There is
// no GSI backing this: legacy tags saved before Scope existed have no scope
// attribute and never match. See docs/adr/0002-scope-filter-without-gsi.md.
func (r *TagDynamoDBRepository) FindTagsByScope(ctx context.Context, scope string) ([]domain.Tag, error) {
	user, err := security.GetCurrentUser(ctx)
	if err != nil {
		return nil, err
	}
	input := &dynamodb.QueryInput{
		TableName: aws.String(r.TableName),
		ExpressionAttributeNames: map[string]string{
			"#scope": "scope", // "scope" is a DynamoDB reserved keyword
		},
		ExpressionAttributeValues: map[string]types.AttributeValue{
			":username": &types.AttributeValueMemberS{Value: *user.UserName},
			":scope":    &types.AttributeValueMemberS{Value: scope},
		},
		KeyConditionExpression: aws.String("user_name = :username"),
		FilterExpression:       aws.String("#scope = :scope"),
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
			Scope: scopeFrom(item),
		})
	}

	return tags, nil
}

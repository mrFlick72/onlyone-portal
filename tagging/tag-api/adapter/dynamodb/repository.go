package dynamodb

import (
	"context"

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
			// scope is always written, defaulting to "" when the caller
			// doesn't supply one. "Optional" describes what the caller must
			// send, not whether the attribute exists in storage — see
			// docs/adr/0003-scope-always-persisted-empty-string-default.md.
			"scope": &types.AttributeValueMemberS{Value: domain.NormalizeScope(tag.Scope)},
		},
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

// FindAllTags queries by user_name. An empty scope means "no filter, return
// everything" rather than "match an empty string" — the only way a tag
// saved before Scope existed (no scope attribute at all) and a
// freshly-saved unscoped tag (scope == "") both stay reachable without a
// backfill.
//
// A non-empty scope is an INCLUSIVE filter: it returns tags whose normalized
// scope matches, PLUS unscoped tags (scope == "" or no scope attribute at
// all). Unscoped tags are shared and must remain reachable from any scoped
// caller without a backfill — the frontend still creates tags through the
// unscoped path, and pre-Scope tags have no scope attribute. Only tags scoped
// to a *different* value (e.g. "revenue" when "expense" is requested) are
// excluded. There is no GSI backing this query. See
// docs/adr/0002-scope-filter-without-gsi.md,
// docs/adr/0003-scope-always-persisted-empty-string-default.md,
// docs/adr/0004-single-find-all-tags-method-with-scope-parameter.md, and
// docs/adr/0006-scoped-query-includes-unscoped-tags.md.
func (r *TagDynamoDBRepository) FindAllTags(ctx context.Context, scope string) ([]domain.Tag, error) {
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
	if scope != "" {
		input.ExpressionAttributeNames = map[string]string{
			"#scope": "scope", // "scope" is a DynamoDB reserved keyword
		}
		input.ExpressionAttributeValues[":scope"] = &types.AttributeValueMemberS{Value: scope}
		input.ExpressionAttributeValues[":empty"] = &types.AttributeValueMemberS{Value: ""}
		// Match the requested scope OR unscoped tags: scope == "" (saved after
		// the Scope field existed) or no scope attribute at all (legacy).
		input.FilterExpression = aws.String("#scope = :scope OR #scope = :empty OR attribute_not_exists(#scope)")
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

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

// FindAllTags queries by user_name and STRICTLY filters by scope: it returns
// only tags whose normalized scope equals the requested scope. There is no
// unfiltered read path and no inclusion of unscoped/legacy or foreign-scope
// tags — scope is authoritative. Callers always pass a real, non-empty scope
// (the only route is GET /api/tags/scope/:scope); pre-Scope and unscoped tags
// are backfilled to a real scope out of band. There is no GSI backing this
// query. See docs/adr/0002-scope-filter-without-gsi.md,
// docs/adr/0004-single-find-all-tags-method-with-scope-parameter.md, and
// docs/adr/0007-scope-mandatory-and-scoped-reads-are-strict.md (which
// superseded 0003 and 0006).
func (r *TagDynamoDBRepository) FindAllTags(ctx context.Context, scope string) ([]domain.Tag, error) {
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

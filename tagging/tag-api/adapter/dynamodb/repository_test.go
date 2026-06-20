package dynamodb

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
	"github.com/mrflick72/onlyone-portal/core-services/golang-web-framework/middleware/security"
)

var TableName = "TestTagsTable"
var client, _ = newDynamoDBClient()

func newStubbedContext() context.Context {
	ctx := context.Background()
	userName := security.UserName("testuser")
	user := security.User{UserName: &userName}
	newCtx := context.WithValue(ctx, "user", user)

	return newCtx
}

func newDynamoDBClient() (*dynamodb.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		config.WithRegion("eu-central-1"),
		config.WithBaseEndpoint("http://localhost:4566"),
	)

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	return dynamodb.NewFromConfig(cfg), err
}

func newTagDynamoDBRepository() *TagDynamoDBRepository {
	return &TagDynamoDBRepository{
		// Initialize with mock or test dependencies as needed
		TableName: TableName,
		Client:    client,
	}
}

func setupTestDynamoDBTable() error {
	// it is an attempt to clean up possible dirty state before creating
	teardownTestDynamoDBTable()
	_, err := client.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: aws.String(TableName),
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("user_name"),
				AttributeType: types.ScalarAttributeTypeS,
			},
			{
				AttributeName: aws.String("search_tag_key"),
				AttributeType: types.ScalarAttributeTypeS,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("user_name"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("search_tag_key"),
				KeyType:       types.KeyTypeRange,
			},
		},
		BillingMode: types.BillingModePayPerRequest,
	})
	if err != nil {

		var resourceInUseException *types.ResourceInUseException
		if !errors.As(err, &resourceInUseException) {
			return err
		}
	}
	return nil
}

func teardownTestDynamoDBTable() error {
	client, _ := newDynamoDBClient()
	_, err := client.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
		TableName: aws.String(TableName),
	})
	return err
}

func TestMain(m *testing.M) {
	setupTestDynamoDBTable()

	code := m.Run() // run all tests

	teardownTestDynamoDBTable()
	os.Exit(code)
}

func findTagByKey(tags []domain.Tag, key string) *domain.Tag {
	for _, tag := range tags {
		if tag.Key == key {
			return &tag
		}
	}
	return nil
}

func saveATag(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "exampleKey", Value: "exampleValue"}

	err := repo.SaveTag(newStubbedContext(), &tag)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFindAllTags(t *testing.T) {
	saveATag(t)

	repo := newTagDynamoDBRepository()

	tags, err := repo.FindAllTags(newStubbedContext(), "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(tags) == 0 {
		t.Errorf("Expected some tags, got none")
	}
}

func TestSaveTagPersistsNormalizedScope(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "scopedKey", Value: "scopedValue", Scope: "  Expense  "}

	if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tags, err := repo.FindAllTags(newStubbedContext(), "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	got := findTagByKey(tags, "scopedKey")
	if got == nil {
		t.Fatalf("Expected to find tag with key %q, got %+v", "scopedKey", tags)
	}
	if got.Scope != "expense" {
		t.Errorf("Expected normalized scope %q, got %q", "expense", got.Scope)
	}
}

func TestSaveTagWithoutScopePersistsEmptyStringScopeAttribute(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "unscopedKey", Value: "unscopedValue"}

	if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tags, err := repo.FindAllTags(newStubbedContext(), "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	got := findTagByKey(tags, "unscopedKey")
	if got == nil {
		t.Fatalf("Expected to find tag with key %q, got %+v", "unscopedKey", tags)
	}
	if got.Scope != "" {
		t.Errorf("Expected empty scope for unscoped tag, got %q", got.Scope)
	}

	// the attribute must actually be present and equal to "" (not absent) —
	// scope is always written from this iteration forward, never omitted.
	rawItem, err := client.GetItem(context.Background(), &dynamodb.GetItemInput{
		TableName: aws.String(TableName),
		Key: map[string]types.AttributeValue{
			"user_name":      &types.AttributeValueMemberS{Value: "testuser"},
			"search_tag_key": &types.AttributeValueMemberS{Value: "unscopedKey"},
		},
	})
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	scopeAttr, ok := rawItem.Item["scope"].(*types.AttributeValueMemberS)
	if !ok {
		t.Fatalf("Expected scope attribute to be present as a string, got %#v", rawItem.Item["scope"])
	}
	if scopeAttr.Value != "" {
		t.Errorf("Expected stored scope attribute to be empty string, got %q", scopeAttr.Value)
	}
}

func TestFindAllTagsWithEmptyScopeReturnsEverythingUnfiltered(t *testing.T) {
	repo := newTagDynamoDBRepository()

	scoped := domain.Tag{Key: "emptyScopeQueryScopedKey", Value: "scopedValue", Scope: "expense"}
	unscoped := domain.Tag{Key: "emptyScopeQueryUnscopedKey", Value: "unscopedValue"}

	for _, tag := range []domain.Tag{scoped, unscoped} {
		if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	tags, err := repo.FindAllTags(newStubbedContext(), "")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	found := map[string]bool{}
	for _, tag := range tags {
		found[tag.Key] = true
	}
	if !found["emptyScopeQueryScopedKey"] {
		t.Errorf("Expected scoped tag to be included in unfiltered (empty-scope) query, got %+v", tags)
	}
	if !found["emptyScopeQueryUnscopedKey"] {
		t.Errorf("Expected unscoped tag to be included in unfiltered (empty-scope) query, got %+v", tags)
	}
}

func TestFindAllTagsWithScopeReturnsOnlyMatchingAndExcludesScopeless(t *testing.T) {
	repo := newTagDynamoDBRepository()

	matching := domain.Tag{Key: "matchKey", Value: "matchValue", Scope: "Revenue"}
	other := domain.Tag{Key: "otherKey", Value: "otherValue", Scope: "expense"}
	legacy := domain.Tag{Key: "legacyScopeKey", Value: "legacyScopeValue"}

	for _, tag := range []domain.Tag{matching, other, legacy} {
		if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	tags, err := repo.FindAllTags(newStubbedContext(), "revenue")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if len(tags) != 1 {
		t.Fatalf("Expected exactly 1 matching tag, got %d: %+v", len(tags), tags)
	}
	if tags[0].Key != "matchKey" {
		t.Errorf("Expected matchKey, got %v", tags[0].Key)
	}
}
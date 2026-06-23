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
	tag := domain.Tag{Key: "exampleKey", Value: "exampleValue", Scope: "expense"}

	err := repo.SaveTag(newStubbedContext(), &tag)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestFindAllTags(t *testing.T) {
	saveATag(t)

	repo := newTagDynamoDBRepository()

	tags, err := repo.FindAllTags(newStubbedContext(), "expense")
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

	tags, err := repo.FindAllTags(newStubbedContext(), "expense")
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

func TestFindAllTagsWithScopeReturnsOnlyExactlyMatchingScope(t *testing.T) {
	repo := newTagDynamoDBRepository()

	matching := domain.Tag{Key: "matchKey", Value: "matchValue", Scope: "Revenue"}
	other := domain.Tag{Key: "otherKey", Value: "otherValue", Scope: "expense"}

	for _, tag := range []domain.Tag{matching, other} {
		if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	}

	tags, err := repo.FindAllTags(newStubbedContext(), "revenue")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	found := map[string]bool{}
	for _, tag := range tags {
		found[tag.Key] = true
	}

	// the requested scope matches (normalized: "Revenue" stored == "revenue" queried)
	if !found["matchKey"] {
		t.Errorf("Expected the scope-matching tag to be included, got %+v", tags)
	}
	// a tag scoped to a *different* value is excluded — strict scoping, no
	// inclusion of foreign scopes or unscoped tags (see ADR 0007).
	if found["otherKey"] {
		t.Errorf("Expected the differently-scoped tag to be excluded, got %+v", tags)
	}
}

func TestUpdateTagValueChangesValueWhenScopeMatches(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "updateKey", Value: "before", Scope: "expense"}
	if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := repo.UpdateTagValue(newStubbedContext(), &domain.Tag{Key: "updateKey", Scope: "expense", Value: "after"}); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tags, err := repo.FindAllTags(newStubbedContext(), "expense")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	got := findTagByKey(tags, "updateKey")
	if got == nil {
		t.Fatalf("Expected to find tag with key %q, got %+v", "updateKey", tags)
	}
	if got.Value != "after" {
		t.Errorf("Expected value %q, got %q", "after", got.Value)
	}
}

func TestUpdateTagValueGivenWrongScopeReturnsNotFound(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "wrongScopeUpdateKey", Value: "before", Scope: "expense"}
	if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err := repo.UpdateTagValue(newStubbedContext(), &domain.Tag{Key: "wrongScopeUpdateKey", Scope: "revenue", Value: "after"})
	if !errors.Is(err, domain.ErrTagNotFound) {
		t.Errorf("Expected ErrTagNotFound, got %v", err)
	}
}

func TestUpdateTagValueGivenMissingKeyReturnsNotFound(t *testing.T) {
	repo := newTagDynamoDBRepository()

	err := repo.UpdateTagValue(newStubbedContext(), &domain.Tag{Key: "missingUpdateKey", Scope: "expense", Value: "after"})
	if !errors.Is(err, domain.ErrTagNotFound) {
		t.Errorf("Expected ErrTagNotFound, got %v", err)
	}
}

func TestDeleteTagRemovesItWhenScopeMatches(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "deleteKey", Value: "value", Scope: "expense"}
	if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if err := repo.DeleteTag(newStubbedContext(), "deleteKey", "expense"); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	tags, err := repo.FindAllTags(newStubbedContext(), "expense")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if findTagByKey(tags, "deleteKey") != nil {
		t.Errorf("Expected the tag to be deleted, got %+v", tags)
	}
}

func TestDeleteTagGivenWrongScopeReturnsNotFound(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "wrongScopeDeleteKey", Value: "value", Scope: "expense"}
	if err := repo.SaveTag(newStubbedContext(), &tag); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	err := repo.DeleteTag(newStubbedContext(), "wrongScopeDeleteKey", "revenue")
	if !errors.Is(err, domain.ErrTagNotFound) {
		t.Errorf("Expected ErrTagNotFound, got %v", err)
	}

	tags, findErr := repo.FindAllTags(newStubbedContext(), "expense")
	if findErr != nil {
		t.Errorf("Expected no error, got %v", findErr)
	}
	if findTagByKey(tags, "wrongScopeDeleteKey") == nil {
		t.Errorf("Expected the tag to remain when the scope didn't match, got %+v", tags)
	}
}

func TestDeleteTagGivenMissingKeyReturnsNotFound(t *testing.T) {
	repo := newTagDynamoDBRepository()

	err := repo.DeleteTag(newStubbedContext(), "missingDeleteKey", "expense")
	if !errors.Is(err, domain.ErrTagNotFound) {
		t.Errorf("Expected ErrTagNotFound, got %v", err)
	}
}
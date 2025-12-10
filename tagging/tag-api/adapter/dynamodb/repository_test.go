package dynamodb

import (
	"context"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)

func newStubbedContext() *context.Context {
	ctx := context.Background()
	userName := domain.UserName("testuser")
	user := &domain.User{UserName: &userName}
	newCtx, _ := domain.SetCurrentUser(user, &ctx)
	return newCtx
}

func newTagDynamoDBRepository() *TagDynamoDBRepository {
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider("xxx", "xxx", "xxx")),
		config.WithRegion("eu-central-1"),
		config.WithBaseEndpoint("http://localhost:8000"),
	)

	if err != nil {
		panic("unable to load SDK config, " + err.Error())
	}

	return &TagDynamoDBRepository{
		// Initialize with mock or test dependencies as needed
		TableName: "TestTagsTable",
		Client:    dynamodb.NewFromConfig(cfg),
	}
}
func TestSaveTag(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "exampleKey", Value: "exampleValue"}

	err := repo.SaveTag(newStubbedContext(), &tag)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetTagBy(t *testing.T) {
	repo := newTagDynamoDBRepository()
	key := "exampleKey"

	tag, err := repo.GetTagBy(newStubbedContext(), key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tag.Key != key {
		t.Errorf("Expected key %v, got %v", key, tag.Key)
	}
}

func TestFindAllTags(t *testing.T) {
	repo := newTagDynamoDBRepository()

	tags, err := repo.FindAllTags(newStubbedContext())
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(*tags) == 0 {
		t.Errorf("Expected some tags, got none")
	}
}

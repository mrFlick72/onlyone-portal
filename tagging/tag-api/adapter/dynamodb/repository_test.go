package dynamodb

import (
	"testing"

	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)
func newTagDynamoDBRepository() *TagDynamoDBRepository {
	return &TagDynamoDBRepository{
		// Initialize with mock or test dependencies as needed	
		TableName: "TestTagsTable",
		UserRepository: nil,
		Client:         nil,
	}
}
func TestSaveTag(t *testing.T) {
	repo := newTagDynamoDBRepository()
	tag := domain.Tag{Key: "exampleKey", Value: "exampleValue"}

	err := repo.SaveTag(tag)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestGetTagBy(t *testing.T) {
	repo := newTagDynamoDBRepository()
	key := "exampleKey"

	tag, err := repo.GetTagBy(key)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if tag.Key != key {
		t.Errorf("Expected key %v, got %v", key, tag.Key)
	}
}

func TestFindAllTags(t *testing.T) {
	repo := newTagDynamoDBRepository()

	tags, err := repo.FindAllTags()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if len(*tags) == 0 {
		t.Errorf("Expected some tags, got none")
	}
}

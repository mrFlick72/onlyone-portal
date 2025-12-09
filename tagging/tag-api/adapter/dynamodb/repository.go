package dynamodb

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/mrFlick72/onlyone-portal/tagging/tag-api/domain"
)


type TagDynamoDBRepository struct {
	UserRepository domain.UserRepository
	Client         *dynamodb.Client
	TableName      string
}

// Implement TagRepository methods here
// For example: SaveTag, GetTagBy, FindAllTags
// These methods will interact with DynamoDB to perform the required operations
// based on the TagRepository interface defined in the domain package.
// You will need to import the domain package to use the Tag struct and TagRepository interface
// import "tagging/tag-api/domain"

func (r *TagDynamoDBRepository) SaveTag(tag domain.Tag) error {
	// Implementation for saving a tag to DynamoDB
	user, err := r.UserRepository.GetCurrentUser()
	if err != nil {
		return err
	}

	input := &dynamodb.PutItemInput{
		TableName: aws.String(r.TableName),
		Item: map[string]types.AttributeValue{
			"Key":      &types.AttributeValueMemberS{Value: tag.Key},
			"Value":    &types.AttributeValueMemberS{Value: tag.Value},
			"UserName": &types.AttributeValueMemberS{Value: *user.UserName},
		},
	}
	_, err = r.Client.PutItem(context.TODO(), input)

	return err
}

func (r *TagDynamoDBRepository) GetTagBy(key string) (*domain.Tag, error) {
	// Implementation for retrieving a tag by key from DynamoDB
	user, err := r.UserRepository.GetCurrentUser()
	if err != nil {
		return nil, err
	}

	input := &dynamodb.GetItemInput{
		TableName: aws.String(r.TableName),
		Key: map[string]types.AttributeValue{
			"Key":      &types.AttributeValueMemberS{Value: key},
			"UserName": &types.AttributeValueMemberS{Value: *user.UserName},
		},
	}
	result, err := r.Client.GetItem(context.TODO(), input)
	if err != nil {
		return nil, err
	}

	if result.Item == nil {
		return nil, nil // Tag not found
	}
	var tag domain.Tag
	err = attributevalue.UnmarshalMap(result.Item, &tag)
	if err != nil {
		return nil, err
	}

	// Return the tag found
	return &tag, nil
}

func (r *TagDynamoDBRepository) FindAllTags() (*[]domain.Tag, error) {
	// Implementation for retrieving all tags from DynamoDB

	return &[]domain.Tag{}, nil
}

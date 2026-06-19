package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type findAllTagsMockRepo struct {
	tags []Tag
	err  error
}

func (m *findAllTagsMockRepo) SaveTag(ctx context.Context, tag *Tag) error { return nil }

func (m *findAllTagsMockRepo) GetTagBy(ctx context.Context, key string) (*Tag, error) { return nil, nil }

func (m *findAllTagsMockRepo) FindAllTags(ctx context.Context) ([]Tag, error) { return m.tags, m.err }

func TestFindAllTagsAppendsUnknownSentinel(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{tags: []Tag{{Key: "tag1", Value: "value1"}}}}

	result, err := action.Execute(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, []Tag{{Key: "tag1", Value: "value1"}, {Key: "UNKNOWN", Value: "UNKNOWN"}}, result)
}

func TestFindAllTagsAppendsUnknownSentinelWhenEmpty(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{tags: []Tag{}}}

	result, err := action.Execute(context.Background())

	assert.NoError(t, err)
	assert.Equal(t, []Tag{{Key: "UNKNOWN", Value: "UNKNOWN"}}, result)
}

func TestFindAllTagsPropagatesRepositoryError(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{err: errors.New("boom")}}

	result, err := action.Execute(context.Background())

	assert.Error(t, err)
	assert.Nil(t, result)
}

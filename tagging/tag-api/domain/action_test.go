package domain

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

type findAllTagsMockRepo struct {
	tags          []Tag
	err           error
	receivedScope string
}

func (m *findAllTagsMockRepo) SaveTag(ctx context.Context, tag *Tag) error { return nil }

func (m *findAllTagsMockRepo) FindAllTags(ctx context.Context, scope string) ([]Tag, error) {
	m.receivedScope = scope
	return m.tags, m.err
}

func TestFindAllTagsAppendsUnknownSentinel(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{tags: []Tag{{Key: "tag1", Value: "value1"}}}}

	result, err := action.Execute(context.Background(), "")

	assert.NoError(t, err)
	assert.Equal(t, []Tag{{Key: "tag1", Value: "value1"}, {Key: "UNKNOWN", Value: "UNKNOWN"}}, result)
}

func TestFindAllTagsAppendsUnknownSentinelWhenEmpty(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{tags: []Tag{}}}

	result, err := action.Execute(context.Background(), "")

	assert.NoError(t, err)
	assert.Equal(t, []Tag{{Key: "UNKNOWN", Value: "UNKNOWN"}}, result)
}

func TestFindAllTagsPropagatesRepositoryError(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{err: errors.New("boom")}}

	result, err := action.Execute(context.Background(), "")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestFindAllTagsWithScopeAppendsUnknownSentinel(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{tags: []Tag{{Key: "tag1", Value: "value1", Scope: "expense"}}}}

	result, err := action.Execute(context.Background(), "expense")

	assert.NoError(t, err)
	assert.Equal(t, []Tag{{Key: "tag1", Value: "value1", Scope: "expense"}, {Key: "UNKNOWN", Value: "UNKNOWN"}}, result)
}

func TestFindAllTagsNormalizesScopeBeforeQuerying(t *testing.T) {
	repo := &findAllTagsMockRepo{}
	action := &FindAllTags{Repository: repo}

	_, err := action.Execute(context.Background(), "  Expense  ")

	assert.NoError(t, err)
	assert.Equal(t, "expense", repo.receivedScope)
}

func TestFindAllTagsWithScopePropagatesRepositoryError(t *testing.T) {
	action := &FindAllTags{Repository: &findAllTagsMockRepo{err: errors.New("boom")}}

	result, err := action.Execute(context.Background(), "expense")

	assert.Error(t, err)
	assert.Nil(t, result)
}

func TestNormalizeScope(t *testing.T) {
	assert.Equal(t, "expense", NormalizeScope("  Expense  "))
	assert.Equal(t, "", NormalizeScope(""))
	assert.Equal(t, "", NormalizeScope("   "))
}

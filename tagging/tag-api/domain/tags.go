package domain

import (
	"context"
	"strings"
)

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Scope string `json:"scope,omitempty"`
}

// NormalizeScope trims and lower-cases a Scope so that values differing
// only by case or surrounding whitespace are treated as the same Scope.
func NormalizeScope(scope string) string {
	return strings.ToLower(strings.TrimSpace(scope))
}

type TagRepository interface {
	SaveTag(ctx context.Context, tag *Tag) error
	GetTagBy(ctx context.Context, key string) (*Tag, error)
	// FindAllTags returns the user's tags. An empty scope returns every tag
	// unfiltered; a non-empty scope returns only tags whose normalized Scope
	// matches it.
	FindAllTags(ctx context.Context, scope string) ([]Tag, error)
}

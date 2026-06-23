package domain

import "context"

type FindAllTags struct {
	Repository TagRepository
}

// Execute returns the user's tags plus the UNKNOWN sentinel. An empty scope
// returns every tag unfiltered; a non-empty scope returns tags whose
// normalized Scope matches it plus unscoped tags (scope "" or unset), so
// shared/legacy tags stay reachable from any scoped caller.
func (action *FindAllTags) Execute(ctx context.Context, scope string) ([]Tag, error) {
	tags, err := action.Repository.FindAllTags(ctx, NormalizeScope(scope))
	if err != nil {
		return nil, err
	}
	return append(tags, Tag{Key: UnknownSentinelKey, Value: UnknownSentinelKey}), nil
}

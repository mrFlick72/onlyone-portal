package tags

import "context"

type SearchTagKey = string
type SearchTagValue = string

type SearchTag struct {
	Key   SearchTagKey   `json:"key"`
	Value SearchTagValue `json:"value"`
}

// UnknownSentinel is the catch-all tag budget-api assigns to a budget record
// (expense or revenue) saved with no tags, so every record resolves to at
// least one tag. It mirrors the UNKNOWN sentinel that tag-api synthesizes into
// every read (tag-api ADR 0001); the two services duplicate the literal with
// no shared enforcement. Defined once here so both aggregates share it.
func UnknownSentinel() SearchTag {
	return SearchTag{Key: "UNKNOWN", Value: "UNKNOWN"}
}

type SearchTagRepository interface {
	GetTagBy(ctx context.Context, key string) (*SearchTag, error)
	GetAllTags(ctx context.Context) ([]SearchTag, error)
}

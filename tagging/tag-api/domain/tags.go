package domain

import (
	"context"
	"errors"
	"strings"
)

// UnknownSentinelKey is the Key of the synthesized UNKNOWN sentinel tag (see
// FindAllTags). It is never persisted, so update/delete reject it outright.
const UnknownSentinelKey = "UNKNOWN"

// ErrTagNotFound is returned by UpdateTagValue/DeleteTag when no tag matches
// the given (key, scope) composite identity — either the key doesn't exist
// at all, or it exists under a different Scope.
var ErrTagNotFound = errors.New("tag not found")

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
	// FindAllTags returns only the tags whose normalized Scope matches scope
	// exactly — Scope is authoritative, see
	// docs/adr/0007-scope-mandatory-and-scoped-reads-are-strict.md.
	FindAllTags(ctx context.Context, scope string) ([]Tag, error)
	// UpdateTagValue changes the Value of the tag identified by tag.Key and
	// tag.Scope — those two fields are the lookup identity, not new values to
	// write; Key and Scope are immutable — see
	// docs/adr/0008-tag-update-is-value-only.md. Returns ErrTagNotFound when
	// no tag matches that composite identity.
	UpdateTagValue(ctx context.Context, tag *Tag) error
	// DeleteTag deletes the tag identified by (key, scope). Returns
	// ErrTagNotFound when no tag matches that composite identity.
	DeleteTag(ctx context.Context, key string, scope string) error
}

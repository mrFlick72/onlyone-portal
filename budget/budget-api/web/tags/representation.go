package tags

import "github.com/mrflick72/budget/budget-api/domain/tags"

type SearchTagRepresentation struct {
	Key   tags.SearchTagKey `json:"tagKey"`
	Value tags.SearchTagValue `json:"tagValue"`
}

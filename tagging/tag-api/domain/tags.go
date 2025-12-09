package domain

type Tag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type TagRepository interface {
	SaveTag(tag Tag) error
	GetTagBy(key string) (*Tag, error)
	FindAllTags() (*[]Tag, error)
}

package cache

type CacheProvider interface {
	Get(key string) (string, bool)
	Set(key string, value string) error
	Evict(key string) error
}

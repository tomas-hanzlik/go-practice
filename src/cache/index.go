package cache

type Cache struct {
	items map[string]string
}

func (cache *Cache) Size() int64 {
	return int64(len(cache.items))
}

func NewCache() *Cache {
	cacheItems := make(map[string]string)

	return &Cache{cacheItems}
}
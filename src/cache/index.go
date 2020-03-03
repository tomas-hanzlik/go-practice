package cache

import (
	// "errors"
	types "./types"
)

type Cache struct {
	items map[string]types.CacheItem
}

func (cache *Cache) Size() int64 {
	return int64(len(cache.items))
}

func (cache *Cache) AddItem(item types.CacheItem) {
	return
}


func NewCache() *Cache {
	cacheItems := make(map[string]types.CacheItem)

	return &Cache{cacheItems}
}

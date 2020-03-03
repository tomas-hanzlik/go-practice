package cache

import (
	"errors"
	types "./types"
)

type Cache struct {
	items map[string]types.CacheItem
}

func (cache *Cache) Size() int64 {
	return int64(len(cache.items))
}

func (cache *Cache) AddItem(item types.CacheItem) error {
	// TODO: Check if cache is full and if yes -> return error
	if false {
		return errors.New("Cache is full.")
	}
	cache.items[item.Key] = item
	return nil
}

func (cache *Cache) GetItem(key string) (types.CacheItem, bool) {
	item, found := cache.items[key]
	
	// TODO: expired check

	return item, found
}


func NewCache() *Cache {
	cacheItems := make(map[string]types.CacheItem)

	return &Cache{cacheItems}
}

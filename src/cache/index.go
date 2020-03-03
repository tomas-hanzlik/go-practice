package cache

import (
	"sync"
	"errors"
	types "./types"
)

// Use RWMutex instead of Mutex for better performance
// with read ops and for safe concurency with `map`
type Cache struct {
	items map[string]types.CacheItem
	m sync.RWMutex
}

func (cache *Cache) Size() int64 {
	return int64(len(cache.items))
}

func (cache *Cache) AddItem(item types.CacheItem) error {
	cache.m.Lock()
	defer cache.m.Unlock()

	// TODO: Check if cache is full and if yes -> return error
	if false {
		return errors.New("Cache is full.")
	}
	cache.items[item.Key] = item
	return nil
}

func (cache *Cache) GetItem(key string) (types.CacheItem, bool) {
	cache.m.RLock()
	defer cache.m.RUnlock()
	
	item, found := cache.items[key]
	
	// TODO: expired check

	return item, found
}


func NewCache() *Cache {
	cacheItems := make(map[string]types.CacheItem)

	return &Cache{items: cacheItems}
}

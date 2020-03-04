package cache // change to CACHE

import (
	"fmt"
	"sync"
	"time"

	types "./types"
)

// Use RWMutex instead of Mutex for better performance
// with read ops and for safe concurency with `map`
type Cache struct {
	Config types.CacheConfig
	Store  map[string]types.CacheItemWrapper
	M      sync.RWMutex
}

func (cache *Cache) Size() int64 {
	return int64(len(cache.Store))
}

func (cache *Cache) AddItem(item types.CacheItem) {
	cache.M.Lock()
	defer cache.M.Unlock()

	newWrappedItem := types.CacheItemWrapper{
		CacheItem:    item,
		ExpirationAt: time.Now().Unix() + int64(cache.Config.Ttl),
	}
	cache.Store[item.Key] = newWrappedItem

	// remove oldest if cache overflow... can happen just once ...
	if cache.Config.Capacity != 0 && cache.Size() > cache.Config.Capacity {
		oldestKey, oldestTimestamp := newWrappedItem.Key, newWrappedItem.ExpirationAt
		for key, wrappedItem := range cache.Store {
			if wrappedItem.ExpirationAt <= oldestTimestamp {
				oldestKey, oldestTimestamp = key, wrappedItem.ExpirationAt
			}
		}
		delete(cache.Store, oldestKey)
	}
}

func (cache *Cache) GetItem(key string) (types.CacheItem, bool) {
	cache.M.RLock()
	defer cache.M.RUnlock()

	wrappedItem, found := cache.Store[key]
	if wrappedItem.IsExpired() {
		delete(cache.Store, key)
		return types.CacheItem{}, false
	}

	return wrappedItem.ToCacheItem(), found
}

func (cache *Cache) RemoveItem(key string) {
	cache.M.Lock()
	defer cache.M.Unlock()

	delete(cache.Store, key)
}

func (cache *Cache) RemoveExpiredItems() {
	return
}

func (cache *Cache) TriggerExpiredItemsRemoval() {
	ticker := time.NewTicker(time.Duration(cache.Config.ExpCheckFrequency) * time.Second)

	for _ = range ticker.C {
		fmt.Println("TOOD")
		cache.RemoveExpiredItems()
	}
}

func NewCache(config types.CacheConfig) *Cache {
	cacheItems := make(map[string]types.CacheItemWrapper)

	cache := &Cache{
		Store:  cacheItems,
		Config: config,
	}
	if config.ExpCheckFrequency > 0 {
		go cache.TriggerExpiredItemsRemoval()
	}

	return cache
}

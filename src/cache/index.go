package cache

import (
	"fmt"
	"os"
	"sync"
	"time"

	types "./types"
)

// Use RWMutex instead of Mutex for better performance
// with read ops and for safe concurency with `map`
type Cache struct {
	Config        types.CacheConfig
	Store         map[string]types.CacheItemWrapper
	InputAdapters []IAdapter
	m             sync.RWMutex
	wg            sync.WaitGroup // to allow input adapters block execution of periodic tasks
}

func (cache *Cache) SetInputAdapter(adapter IAdapter) {
	// Wait for needy adapter if needed.
	cache.InputAdapters = append(cache.InputAdapters, adapter)
	cache.InputAdapters[len(cache.InputAdapters)-1].Run(&cache.wg) // TODO: pointers
}

func (cache *Cache) CollectAdaptersData() {
	for _, adapter := range cache.InputAdapters {
		for _, item := range adapter.GetData() {
			cache.AddItem(*item)
		}

		a, ok := adapter.(INoisyAdapter)
		if ok {
			fmt.Println(a.Stats())
		}
	}
	return
}

func (cache *Cache) Size() int64 {
	return int64(len(cache.Store))
}

func (cache *Cache) AddItem(item types.CacheItem) {
	cache.m.Lock()
	defer cache.m.Unlock()

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
	cache.m.RLock()
	defer cache.m.RUnlock()

	wrappedItem, found := cache.Store[key]
	if wrappedItem.IsExpired() {
		delete(cache.Store, key)
		return types.CacheItem{}, false
	}

	return wrappedItem.ToCacheItem(), found
}

func (cache *Cache) RemoveItem(key string) {
	cache.m.Lock()
	defer cache.m.Unlock()

	delete(cache.Store, key)
}

func (cache *Cache) RemoveExpiredItems() {
	for key, wrappedItem := range cache.Store {
		if wrappedItem.IsExpired() {
			cache.RemoveItem(key)
		}
	}
}

func (cache *Cache) Dump(filename string) {
	file, err := os.Create(filename)

	if err != nil {
		panic(err)
	}
	defer file.Close()

	for _, item := range cache.Store {
		file.WriteString(fmt.Sprintf("%s:%s\n", item.Key, item.Value))
	}
}

func NewCache(config types.CacheConfig) *Cache {
	cacheItems := make(map[string]types.CacheItemWrapper, 0)
	cache := &Cache{
		Store:  cacheItems,
		Config: config,
	}

	// Collect data from adapters.
	if cache.Config.GetDataFrequency > 0 
		ExecutePeriodic(&cache.wg, cache.Config.GetDataFrequency, cache.CollectAdaptersData)
	}

	// Remove expired items.
	if cache.Config.ExpCheckFrequency > 0 {
		ExecutePeriodic(&cache.wg, cache.Config.ExpCheckFrequency, cache.RemoveExpiredItems)
	}
	return cache
}

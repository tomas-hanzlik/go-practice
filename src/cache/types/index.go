package types

import (
	"time"
)

type CacheItem struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// wrap cache item for internal usage of cache manager
type CacheItemWrapper struct {
	CacheItem
	ExpirationAt int64
}

func (item *CacheItemWrapper) IsExpired() bool {
	return item.ExpirationAt <= time.Now().Unix()
}

func (item *CacheItemWrapper) ToCacheItem() CacheItem {
	return CacheItem{
		Key:   item.Key,
		Value: item.Value,
	}
}

type CacheConfig struct {
	TTL                      int32 `json:"ttl"`                      // Expiration of items.
	Capacity                 int64 `json:"capacity"`                 // Capacity of the cache.
	ExpCheckFrequency        int32 `json:"expirationCheckFrequency"` // How often remove expired items. 0 to turn it off
	GetAdaptersDataFrequency int32 `json:"getAdaptersDataFrequency"` // How often we want to get data from adapters
	AdaptersBufferSize       int64 `json:"adaptersBufferSize"`       // If we want to limit the amount of data before colleciton
}

// Use as a custom buffer
type ItemsQueue struct {
	items    []CacheItem // items queue
	Capacity int64       // 0 for unlimited
}

func (q *ItemsQueue) Size() int64 {
	return int64(len(q.items))
}

func (q *ItemsQueue) IsEmpty() bool {
	return q.Size() == 0
}

func (q *ItemsQueue) Enq(item CacheItem) {
	q.items = append(q.items, item)

	if q.Capacity > 0 && q.Size() > q.Capacity {
		q.Deq()
	}
}

func (q *ItemsQueue) Deq() CacheItem {
	if q.IsEmpty() {
		return CacheItem{}
	}

	item := q.items[0]
	q.items = q.items[1:]
	return item
}

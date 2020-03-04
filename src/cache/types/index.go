package types

import "time"

type CacheItem struct {
	Key   string // To satify the assigment methods desciption (addItem, Adapters, ...)
	Value string
}

// wrap cache item for internal usage of cache manager
type CacheItemWrapper struct {
	CacheItem
	ExpirationAt int64
}

// add method if expired
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
	Ttl               int32 // Expiration of items.
	Capacity          int64 // Capacity of the cache.
	ExpCheckFrequency int32 // How often remove expired items. 0 to turn it off
}

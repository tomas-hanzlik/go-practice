package types

// Add CacheItem type with one string field for storing cache value. Put its definition to src/cache/types dir. Add new method called Ad
// dItem to Cache which would have one argument - of CacheItem type. You will also need to create some container for cache items.
// AddItem should be concurrent, work with sync.Map or implement your own logic.
type CacheItem struct {
	Key   string // To preserve the interface from other tasks (addItem, Adapters, ...)... TOOD: Ask
	Value string
}

type CacheConfig struct {
	Ttl               int32
	Size              int64
	ExpCheckFrequency int64
}

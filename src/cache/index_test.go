package cache

import (
	"fmt"
	"testing"

	types "./types"
	"github.com/stretchr/testify/assert"
)

// test data
var defaultKeys = map[string]string{
	"one":   "tomas",
	"two":   "2138",
	"three": "1908",
	"four":  "912",
}

type MockCache struct {
	*Cache
}

// Helper function to fill cache with data to make testing of some fucntionality easier
func (cache *MockCache) FillWithDefaultData() {
	for key, value := range defaultKeys {
		cache.AddItem(types.CacheItem{Key: key, Value: value})
	}
}

// Helper function to fill cache with expired data.
// nasty implementation ... miss Time mock :'(
func (cache *MockCache) FillWithDefaultDataAsExpired() {
	for key, value := range defaultKeys {
		cache.Store[key] = types.CacheItemWrapper{
			CacheItem: types.CacheItem{
				Key:   key,
				Value: value,
			},
			ExpirationAt: 1,
		}
	}
}

func NewMockCache() *MockCache {
	config := types.CacheConfig{
		Ttl:               30,
		Capacity:          100,
		ExpCheckFrequency: 0, // to disable periodic expiration check
	}
	return &MockCache{Cache: NewCache(config)}
}

// helper function that can abstract out some logic and prepares testing state - in this case it creates basic cache instance
func prepareBrandNewCache() *MockCache {
	cache := NewMockCache()

	// fmt.Println(cache.config)
	return cache
}

// Tests if constructor works
func TestNewCache(t *testing.T) {
	cache := NewMockCache()

	assert.NotNil(t, cache, "cache constructor should not return nil")
}

// unit test for specific method
func TestCache_Size(t *testing.T) {
	cache := prepareBrandNewCache()

	assert.Empty(t, cache.Size(), "brand new cache should have zero items")
}

func TestCache_AddItem(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.Config.Capacity = 1

	cache.AddItem(types.CacheItem{Key: "t1", Value: "42"})
	assert.Equal(t, int64(1), cache.Size(), "Cache should have exactly one item.")

	// check cache overflow and removal of the oldest one
	cache.AddItem(types.CacheItem{Key: "t2", Value: "42"})
	assert.Equal(t, int64(1), cache.Size(), "Cache should have exactly one item.")
	assert.NotNil(t, cache.Store["t2"], "Cache should have exactly one item `t2`.")
}

func TestCache_GetItem(t *testing.T) {
	cache := prepareBrandNewCache()
	// Test getting a known item
	cache.FillWithDefaultData() // helper
	item, found := cache.GetItem("one")
	assert.True(t, found, "In case of known item should return true.")
	assert.Equal(t, item.Value, defaultKeys["one"], "In case of known item should correct item.")

	// Test if we can get expired Item
	cache.FillWithDefaultDataAsExpired()
	item, found = cache.GetItem("one")
	assert.False(t, found, "Expired Item shouldnt be returned.")
	assert.Empty(t, item, "Expired item should be empty.")

	// Test getting unknown item
	item, found = cache.GetItem("UKNOWN_KEY")
	fmt.Println(item)
	assert.False(t, found, "In case of unknown item should return false.")
	assert.Empty(t, item, "In case of unknown item should be empty.")
}

func TestCache_RemoveItem(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.FillWithDefaultData() // helper

	prevSize := cache.Size()

	cache.RemoveItem("one")

	assert.Equal(t, cache.Size(), prevSize-1, "Size of cache after deletion doesnt match.")

}

func TestCache_RemoveExpiredItems(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.FillWithDefaultDataAsExpired()

	cache.RemoveExpiredItems()
	assert.Empty(t, cache.Size(), "Cache should be empty after removal of expired items.")
}

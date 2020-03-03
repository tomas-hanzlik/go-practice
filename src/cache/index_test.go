package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	// "fmt"
	types "./types"
)


type MockCache struct {
	*Cache
}

func NewMockCache() *MockCache {
	return &MockCache{Cache: NewCache()}
}


// Tests if constructor works
func TestNewCache(t *testing.T) {
	cacheInst := NewMockCache()

	assert.NotNil(t, cacheInst, "cache constructor should not return nil")
}

// helper function that can abstract out some logic and prepares testing state - in this case it creates basic cache instance
func prepareBrandNewCache() *MockCache {
	return NewMockCache()
}

// unit test for specific method
func TestCache_Size(t *testing.T) {
	cache := prepareBrandNewCache()

	assert.Empty(t, cache.Size(), "brand new cache should have zero items")
}


// Test adding of new items into the cache
func TestCache_AddItem(t *testing.T) {
	cache := prepareBrandNewCache()

	cache.AddItem(types.CacheItem{"32", "43"})
	assert.Equal(t, int64(1), cache.Size(), "Cache should have exactly one item.")
}
package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	types "./types"
)

// add global list of key during init

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
func TestCache_AddItem_GetItem(t *testing.T) {
	// TODO: Split use cases
	cache := prepareBrandNewCache()

	// test AddItem
	newItem := types.CacheItem{Key: "32", Value: "43"}
	e := cache.AddItem(newItem)
	assert.Equal(t, int64(1), cache.Size(), "Cache should have exactly one item.")
	assert.Nil(t, e, "Failed to add item")

	// test GetItem
	item, found := cache.GetItem("34")
	fmt.Println(item, found)
	assert.False(t, found, "In case of unknown item should return false.")	
	assert.Empty(t, item, "In case of unknown item should be empty.")

	item, found = cache.GetItem("32")
	fmt.Println(item, found)
	assert.True(t, found, "In case of known item should return true.")	
	assert.Equal(t, item, newItem, "In case of known item should correct item.")

}

package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"fmt"
	types "./types"
)

// add global list of key during init


var defaultKeys = map[string]string{
    "one": "tomas",
    "two": "2138",
    "three": "1908",
    "four": "912",
}

type MockCache struct {
	*Cache
}

func (cache *MockCache) FillWithDefaultData() {
	for key, value := range defaultKeys {
		cache.AddItem(types.CacheItem{Key: key, Value: value})
	}
}

func NewMockCache() *MockCache {
	return &MockCache{Cache: NewCache()}
}


// Tests if constructor works
func TestNewCache(t *testing.T) {
	cache := NewMockCache()

	assert.NotNil(t, cache, "cache constructor should not return nil")
}

// helper function that can abstract out some logic and prepares testing state - in this case it creates basic cache instance
func prepareBrandNewCache() *MockCache {
	cache := NewMockCache()

	return cache
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
	// cache.FillWithDefaultData()

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
	fmt.Println(item, found, defaultKeys)
	assert.True(t, found, "In case of known item should return true.")	
	assert.Equal(t, item, newItem, "In case of known item should correct item.")

}


func TestCache_AddItem(t *testing.T) {
	cache := prepareBrandNewCache()

	e := cache.AddItem(types.CacheItem{Key: "32", Value: "43"})
	assert.Equal(t, int64(1), cache.Size(), "Cache should have exactly one item.")
	assert.Nil(t, e, "Failed to add item")
}



func TestCache_GetItem(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.FillWithDefaultData() // helper

	// test GetItem
	item, found := cache.GetItem("SomeRandomKey")
	assert.False(t, found, "In case of unknown item should return false.")	
	assert.Empty(t, item, "In case of unknown item should be empty.")

	item, found = cache.GetItem("one")
	assert.True(t, found, "In case of known item should return true.")	
	assert.Equal(t, item.Value, defaultKeys["one"], "In case of known item should correct item.")
}

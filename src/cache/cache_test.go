package cache

import (
	"fmt"
	"strings"
	"testing"

	types "tohan.net/go-practice/src/cache/types"

	"github.com/kami-zh/go-capturer"
	"github.com/stretchr/testify/assert"
)

// test data
var defaultTestKeyValues = map[string]string{
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
	for key, value := range defaultTestKeyValues {
		cache.AddItem(types.CacheItem{Key: key, Value: value})
	}
}

// Helper function to fill cache with expired data.
func (cache *MockCache) FillWithDefaultDataAsExpired() {
	for key, value := range defaultTestKeyValues {
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
		TTL:                      30,
		Capacity:                 100,
		ExpCheckFrequency:        0, // to disable periodic expiration check
		GetAdaptersDataFrequency: 0,
	}
	return &MockCache{Cache: NewCache(config)}
}

// helper function that can abstract out some logic and prepares testing state - in this case it creates basic cache instance
func prepareBrandNewCache() *MockCache {
	cache := NewMockCache()

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
	assert.Equal(t, item.Value, defaultTestKeyValues["one"], "In case of known item should correct item.")

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

func TestCache_GetAllItems(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.FillWithDefaultData()

	items := cache.GetAllItems()
	assert.Equal(t, len(*items), len(defaultTestKeyValues), "not getting all items")
}

func TestCache_RemoveItem(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.FillWithDefaultData() // helper

	prevSize := cache.Size()

	cache.RemoveItem("one")

	assert.Equal(t, cache.Size(), prevSize-1, "Size of cache after deletion doesnt match.")
}

func TestCache_RemoveAllItems(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.FillWithDefaultData()

	cache.RemoveAllItems()
	assert.Empty(t, cache.Size(), "cache is not empty")
}

func TestCache_RemoveExpiredItems(t *testing.T) {
	cache := prepareBrandNewCache()
	cache.FillWithDefaultDataAsExpired()

	cache.RemoveExpiredItems()
	assert.Empty(t, cache.Size(), "Cache should be empty after removal of expired items.")
}

func TestCache_RandomInputAdapter(t *testing.T) {
	// Random input adapter should be tested separately and then as the whole with `cache`
	// Time constraints... Can be improved if necessary
	cache := prepareBrandNewCache()

	// set data generation frequency to 0 so we can do it manualy
	cache.SetInputAdapter(NewRandomInputAdapter(0, 10, 0))

	// generate data manualy
	cache.InputAdapters[0].(*RandomInputAdapter).generateData()

	// capture StdOut
	out := capturer.CaptureStdout(func() {
		cache.CollectAdaptersData()
	})
	assert.Contains(t, out, "Collecting items:", "Stats() method should be called.")
	assert.NotEmpty(t, cache.Size(), "cache should have randomly generated items")
	assert.True(t, cache.InputAdapters[0].(*RandomInputAdapter).queue.IsEmpty(), "adapter's queue should be empty")

}

func TestCache_CommandLineInputAdapter(t *testing.T) {
	// Random input adapter should be tested separately and then as the whole with `cache`
	// Time constraints... Can be improved if necessary
	testString := `
		test1: test1
		test2: test2
		fsdfsdfs
		test3: test3
		STOP
	`
	cache := prepareBrandNewCache()
	cache.SetInputAdapter(NewCommandLineInputAdapter(strings.NewReader(testString), 0))

	cache.InputAdapters[0].(*CommandLineInputAdapter).readFromStdin()

	cache.CollectAdaptersData()
	assert.Equal(t, int64(3), cache.Size(), "cache size not matching")
	assert.True(t, cache.InputAdapters[0].(*CommandLineInputAdapter).queue.IsEmpty(), "adapter's queue should be empty")
}

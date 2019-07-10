package cache

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

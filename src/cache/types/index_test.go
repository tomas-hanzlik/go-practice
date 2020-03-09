package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestItemsQueue(t *testing.T) {
	queue := ItemsQueue{}
	assert.True(t, queue.IsEmpty(), "queue should be empty")

	queue.Enq(CacheItem{Key: "1", Value: "1"})
	queue.Enq(CacheItem{Key: "2", Value: "2"})
	assert.False(t, queue.IsEmpty(), "queue should not be empty")

	// check if returned item correct
	item := queue.Deq()
	assert.Equal(t, "1", item.Value, "item values should match")

	item = queue.Deq()
	assert.True(t, queue.IsEmpty(), "queue should be empty")

	item = queue.Deq()
	assert.Empty(t, item, "expected empty item")
}

func TestItemsQueue_Capacity(t *testing.T) {
	queue := ItemsQueue{capacity: 2}
	assert.True(t, queue.IsEmpty(), "queue should be empty")

	queue.Enq(CacheItem{Key: "1", Value: "1"})
	queue.Enq(CacheItem{Key: "2", Value: "2"})
	queue.Enq(CacheItem{Key: "3", Value: "3"})

	assert.Equal(t, 2, int(queue.Size()), "Q should be of size 2")

	// check if returned item correct
	item := queue.Deq()
	assert.Equal(t, "2", item.Value, "item values should match")

}

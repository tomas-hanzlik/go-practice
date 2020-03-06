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
	assert.Equal(t, item.Value, "1", "item values should match")

	item = queue.Deq()
	assert.True(t, queue.IsEmpty(), "queue should be empty")

	item = queue.Deq()
	assert.Empty(t, item, "expected empty item")
}

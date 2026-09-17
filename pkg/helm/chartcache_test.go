package helm

import (
	"fmt"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLruChartCache_Add(t *testing.T) {
	t.Run("should store without evicting", func(t *testing.T) {
		sut := newLRUChartCache(2)

		evicted := sut.Add("a", []byte("chart-a"))
		assert.False(t, evicted)
		assert.Len(t, sut.items, 1)
		assert.Equal(t, sut.order.Len(), 1)

		item, ok := sut.items["a"]
		assert.True(t, ok)

		cacheItem, ok := item.Value.(*cacheEntry)
		assert.True(t, ok)
		assert.Equal(t, []byte("chart-a"), cacheItem.value)

		assert.Equal(t, item, sut.order.Front())
	})

	t.Run("should evict the least-recently-used entry when capacity is exceeded", func(t *testing.T) {
		sut := newLRUChartCache(2)
		sut.Add("a", []byte("chart-a"))
		sut.Add("b", []byte("chart-b"))

		// adding a third entry must evict "a" (the least recently used)
		evicted := sut.Add("c", []byte("chart-c"))
		assert.True(t, evicted)

		_, okA := sut.items["a"]
		assert.False(t, okA)
		assert.Len(t, sut.items, 2)
		assert.Equal(t, sut.order.Len(), 2)

		itemB, okB := sut.items["b"]
		assert.True(t, okB)
		itemC, okC := sut.items["c"]
		assert.True(t, okC)

		assert.Equal(t, itemC, sut.order.Front())
		assert.Equal(t, itemB, sut.order.Back())
	})

}

func TestLruChartCache_Get(t *testing.T) {
	t.Run("should return not-ok for missing key", func(t *testing.T) {
		sut := newLRUChartCache(2)

		value, ok := sut.Get("missing")

		assert.False(t, ok)
		assert.Nil(t, value)
	})

	t.Run("should store and return a value", func(t *testing.T) {
		sut := newLRUChartCache(2)

		evicted := sut.Add("a", []byte("chart-a"))

		value, ok := sut.Get("a")
		assert.False(t, evicted)
		assert.True(t, ok)
		assert.Equal(t, []byte("chart-a"), value)
	})

	t.Run("should update the value and recency of an existing key without growing", func(t *testing.T) {
		sut := newLRUChartCache(2)
		sut.Add("a", []byte("old"))

		evicted := sut.Add("a", []byte("new"))

		value, ok := sut.Get("a")
		assert.False(t, evicted)
		assert.True(t, ok)
		assert.Equal(t, []byte("new"), value)
		assert.Equal(t, 1, sut.order.Len())
	})
}

// run with -race to detect data races on the shared cache
func TestLruChartCache_Concurrent(t *testing.T) {
	sut := newLRUChartCache(50)

	var wg sync.WaitGroup
	for i := range 100 {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			key := fmt.Sprintf("chart:%d", n%25)
			sut.Add(key, []byte(key))
			_, _ = sut.Get(key)
		}(i)
	}
	wg.Wait()

	assert.LessOrEqual(t, sut.order.Len(), 50)
	assert.Equal(t, sut.order.Len(), len(sut.items))
}

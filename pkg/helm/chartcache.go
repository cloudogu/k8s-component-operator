package helm

import (
	"container/list"
	"sync"
)

// lruChartCache is a small, thread-safe, fixed-capacity LRU cache for chart archive bytes keyed by "chartName:version".
//
// It exists to avoid pulling the same chart version from the registry on every reconcile (e.g. while a component is
// stuck retrying). A chart:version artifact is immutable by convention, so entries never need time-based invalidation.
// The cache is bounded to avoid increasing memory footprint of the operator. The least-recently-used entry is
// evicted once the capacity is exceeded, which keeps an actively-retried chart resident while idle old versions age out.
type lruChartCache struct {
	mu       sync.Mutex
	capacity int
	// order holds *cacheEntry values with the most-recently-used entry at the front.
	order *list.List
	items map[string]*list.Element
}

type cacheEntry struct {
	key   string
	value []byte
}

// newLRUChartCache creates an LRU cache that holds at most capacity entries. The caller must ensure capacity > 0.
func newLRUChartCache(capacity int) *lruChartCache {
	return &lruChartCache{
		capacity: capacity,
		order:    list.New(),
		items:    make(map[string]*list.Element, capacity),
	}
}

// Get returns the cached bytes for key and whether the key was present, marking the entry as most-recently-used.
func (c *lruChartCache) Get(key string) (value []byte, ok bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	element, ok := c.items[key]
	if !ok {
		return nil, false
	}

	c.order.MoveToFront(element)

	tCacheEntry, _ := element.Value.(*cacheEntry)
	return tCacheEntry.value, true
}

// Add stores value under key as the most-recently-used entry and reports whether an entry was evicted to stay within
// capacity. Adding an existing key updates its value and refreshes its recency.
func (c *lruChartCache) Add(key string, value []byte) (evicted bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, ok := c.items[key]; ok {
		tCacheEntry, _ := element.Value.(*cacheEntry)
		tCacheEntry.value = value
		c.order.MoveToFront(element)

		return false
	}

	c.items[key] = c.order.PushFront(&cacheEntry{key: key, value: value})

	if c.order.Len() > c.capacity {
		c.removeOldest()
		return true
	}

	return false
}

// removeOldest evicts the least-recently-used entry. The caller must hold c.mu.
func (c *lruChartCache) removeOldest() {
	oldest := c.order.Back()
	if oldest == nil {
		return
	}

	c.order.Remove(oldest)
	tCacheEntry, _ := oldest.Value.(*cacheEntry)

	delete(c.items, tCacheEntry.key)
}

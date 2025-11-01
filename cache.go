package main

import (
	"container/list"
	"sync"
	"time"
)

// TTLCache implements a simple LRU cache with TTL
type TTLCache struct {
	mu         sync.RWMutex // Use RWMutex to allow concurrent reads
	maxSize    int
	ttl        time.Duration
	cache      map[string]*list.Element
	lruList    *list.List
	timestamps map[string]time.Time
}

type cacheEntry struct {
	key   string
	value interface{}
}

// NewTTLCache creates a new TTL cache
func NewTTLCache(maxSize int, ttl time.Duration) *TTLCache {
	return &TTLCache{
		maxSize:    maxSize,
		ttl:        ttl,
		cache:      make(map[string]*list.Element),
		lruList:    list.New(),
		timestamps: make(map[string]time.Time),
	}
}

// Get retrieves a value from the cache
func (c *TTLCache) Get(key string) (interface{}, bool) {
	// Use read lock first for fast path
	c.mu.RLock()
	_, exists := c.cache[key]
	if !exists {
		c.mu.RUnlock()
		return nil, false
	}

	// Check TTL (still need to check under read lock)
	timestamp, ok := c.timestamps[key]
	expired := !ok || time.Since(timestamp) > c.ttl
	c.mu.RUnlock()

	// If expired, upgrade to write lock to remove
	if expired {
		c.mu.Lock()
		// Double-check after acquiring write lock
		if element, exists := c.cache[key]; exists {
			timestamp, ok := c.timestamps[key]
			if !ok || time.Since(timestamp) > c.ttl {
				c.lruList.Remove(element)
				delete(c.cache, key)
				delete(c.timestamps, key)
			}
		}
		c.mu.Unlock()
		return nil, false
	}

	// Upgrade to write lock only for LRU update
	c.mu.Lock()
	// Double-check element still exists
	if element, exists := c.cache[key]; exists {
		c.lruList.MoveToBack(element)
		value := element.Value.(*cacheEntry).value
		c.mu.Unlock()
		return value, true
	}
	c.mu.Unlock()
	return nil, false
}

// Put adds a value to the cache
func (c *TTLCache) Put(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if element, exists := c.cache[key]; exists {
		// Update existing item
		c.lruList.MoveToBack(element)
		element.Value.(*cacheEntry).value = value
		c.timestamps[key] = time.Now()
		return
	}

	// Add new item
	if c.lruList.Len() >= c.maxSize {
		// Remove least recently used item
		oldest := c.lruList.Front()
		if oldest != nil {
			entry := oldest.Value.(*cacheEntry)
			c.lruList.Remove(oldest)
			delete(c.cache, entry.key)
			delete(c.timestamps, entry.key)
		}
	}

	entry := &cacheEntry{key: key, value: value}
	element := c.lruList.PushBack(entry)
	c.cache[key] = element
	c.timestamps[key] = time.Now()
}

// Size returns the current size of the cache
func (c *TTLCache) Size() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.lruList.Len()
}

// Clear removes all items from the cache
func (c *TTLCache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache = make(map[string]*list.Element)
	c.lruList = list.New()
	c.timestamps = make(map[string]time.Time)
}


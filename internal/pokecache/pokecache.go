package pokecache

import (
	"sync"
	"time"
)

type Cache struct {
	mu sync.Mutex
	value map[string]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache(interval time.Duration) *Cache {
	c := &Cache{
		value: map[string]cacheEntry{},
		interval: interval,
	}
	go c.reapLoop()
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	newCacheEntry := cacheEntry{
		createdAt: time.Now(),
		val: val,
	}
	c.value[key] = newCacheEntry
	}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	entry, ok := c.value[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mu.Lock()
		for key, entry := range c.value {
			if time.Since(entry.createdAt) >= c.interval {
				delete(c.value, key)
		}
		}
		c.mu.Unlock()
	}
}
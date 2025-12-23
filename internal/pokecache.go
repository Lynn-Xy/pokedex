package pokecache

import (
	"fmt"
	"sync"
	"time"
	"errors"
)

type safeCache struct {
	mu sync.Mutex
	value map[string]cacheEntry
	interval time.Duration
}

type cacheEntry struct {
	createdAt time.Time
	val []byte
}

func NewCache() (safeCache, error) {
	c := safeCache{
		mu: sync.Mutex{}
		value: map[string]cacheEntry{}
	}
}

func (c *safeCache) Add(key string, val []byte) error {
	if len(c.value) == 0 {
		return errors.New("error adding cache entry: cache value is empty")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.value[key]; ok == false {
		mins, _ := time.ParseDuration("5m")
		newCacheEntry := cacheEntry{
			createdAt: time.Now(),
			interval: mins,
			value: val,
		}
		c.value[key] = newCacheEntry
		return nil
	} else {
		return errors.New("error adding cache entry: entry already exists in cache")
	}
}

func (c *safeCache) Get(key string) (cacheEntry, bool, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.value[key]; ok == true {
		return c.value[key], true, nil
	} else {
		return nil, false, errors.New("error getting cache entry: entry does not exist in cache")
	}
}

func (c *safeCache) reapLoop() {
	c.mu.Lock()
	defer c.mu.Unlock()
	for key, entry := range c.value {
		if time.Since(entry.createdAt) >= c.interval {
			delete(c.value, key)
	}
}
}



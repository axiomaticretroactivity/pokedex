package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	entries  map[string]cacheEntry
	interval time.Duration
	mux      sync.RWMutex
}

func (c *Cache) Add(key string, value []byte) {
	c.mux.Lock()
	if c.entries == nil {
		c.entries = map[string]cacheEntry{}
	}
	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       value,
	}
	c.mux.Unlock()
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	entry, ok := c.entries[key]
	if !ok {
		return nil, false
	}
	return entry.val, true
}

func (c *Cache) reap() {
	c.mux.Lock()
	defer c.mux.Unlock()

	now := time.Now()
	interval := c.interval
	cutoff := now.Add(-interval)

	for key, entry := range c.entries {
		if entry.createdAt.Before(cutoff) {
			delete(c.entries, key)
		}
	}
}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.reap()
	}
}

func NewCache(newInterval time.Duration) *Cache {
	newCache := &Cache{
		entries:  map[string]cacheEntry{},
		interval: newInterval,
		mux:      sync.RWMutex{},
	}
	go newCache.reapLoop()
	return newCache
}

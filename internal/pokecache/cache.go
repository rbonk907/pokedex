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
	mu      *sync.Mutex
	entries map[string]cacheEntry
}

func (c *Cache) Add(key string, val []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) (val []byte, found bool) {
	c.mu.Lock()
	defer c.mu.Unlock()

	entry, exists := c.entries[key]
	if exists {
		return entry.val, exists
	}

	return []byte{}, exists
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for t := range ticker.C {
		c.mu.Lock()
		for k, v := range c.entries {
			if t.Sub(v.createdAt) > interval {
				delete(c.entries, k)
			}
		}
		c.mu.Unlock()
	}
}

func NewCache(interval time.Duration) Cache {
	c := Cache{
		mu:      &sync.Mutex{},
		entries: make(map[string]cacheEntry),
	}

	// start the reap loop concurrently
	go c.reapLoop(interval)

	return c
}

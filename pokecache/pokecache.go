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
	entries map[string]cacheEntry
	mu      *sync.Mutex
}

func NewCache(interval time.Duration) Cache {
	cache := Cache{
		make(map[string]cacheEntry),
		&sync.Mutex{},
	}

	go cache.reapLoop(interval)
	return cache

}

func (c Cache) Add(key string, data []byte) error {
	entry := cacheEntry{
		createdAt: time.Now(),
		val:       data,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.entries[key] = entry
	return nil
}

func (c Cache) Get(key string) ([]byte, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	val, found := c.entries[key]
	if found {
		return val.val, found
	}
	return nil, false
}

func (c Cache) reapLoop(interval time.Duration) {

	// create time.Ticker
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		t := <-ticker.C
		c.mu.Lock()
		for key, value := range c.entries {
			if t.Sub(value.createdAt) > interval {
				delete(c.entries, key)
			}
		}

		c.mu.Unlock()
	}

}

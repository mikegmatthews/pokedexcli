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
	lock    sync.Mutex
	entries map[string]cacheEntry
}

func NewCache(maxAge time.Duration) *Cache {
	newCache := Cache{
		lock:    sync.Mutex{},
		entries: make(map[string]cacheEntry),
	}

	reapTicker := time.NewTicker(maxAge)
	go func(c *Cache, tick *time.Ticker) {
		for range tick.C {
			if c == nil {
				tick.Stop()
				return
			}

			c.reapLoop(maxAge)
		}
	}(&newCache, reapTicker)

	return &newCache
}

func (c *Cache) Add(key string, val []byte) {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.entries[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	defer c.lock.Unlock()

	entry, ok := c.entries[key]

	if ok {
		return entry.val, true
	} else {
		return nil, false
	}
}

func (c *Cache) reapLoop(maxAge time.Duration) {
	c.lock.Lock()
	defer c.lock.Unlock()

	for key, entry := range c.entries {
		if entry.createdAt.Add(maxAge).Before(time.Now()) {
			delete(c.entries, key)
		}
	}
}

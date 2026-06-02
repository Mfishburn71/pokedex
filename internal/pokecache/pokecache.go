package pokecache

import (
	"context"
	"sync"
	"time"
)

type cacheEntry struct {
	createdAt time.Time
	val       []byte
}

type Cache struct {
	cacheMap map[string]cacheEntry
	mux      sync.RWMutex
	interval time.Duration
}

func NewCache(interval time.Duration, ctx context.Context) *Cache {
	c := &Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
	go c.reapLoop(ctx)
	return c
}

func (c *Cache) Add(key string, val []byte) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.cacheMap[key] = cacheEntry{
		createdAt: time.Now(),
		val:       val,
	}
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mux.RLock()
	defer c.mux.RUnlock()
	elem, ok := c.cacheMap[key]
	return elem.val, ok
}

func (c *Cache) reapLoop(ctx context.Context) {
	ticker := time.NewTicker(c.interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			c.mux.Lock()
			for key, entry := range c.cacheMap {
				if time.Since(entry.createdAt) > c.interval {
					delete(c.cacheMap, key)
				}
			}
			c.mux.Unlock()
		case <-ctx.Done():
			return
		}
	}
}

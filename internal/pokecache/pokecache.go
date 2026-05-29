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
	cacheMap map[string]cacheEntry
	mux      sync.RWMutex //readWrite Mutex allows multiple readers
	interval time.Duration
}

func NewCache(interval time.Duration) *Cache {

	cacheObj := Cache{
		cacheMap: make(map[string]cacheEntry),
		interval: interval,
	}
	go cacheObj.reapLoop()
	return &cacheObj
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
	if ok == true {
		return elem.val, true
	}
	return nil, false

}

func (c *Cache) reapLoop() {
	ticker := time.NewTicker(c.interval)
	for range ticker.C {
		c.mux.Lock()
		for key, entry := range c.cacheMap {
			if time.Since(entry.createdAt) > c.interval {
				delete(c.cacheMap, key)
			}
		}
		c.mux.Unlock()
	}
}

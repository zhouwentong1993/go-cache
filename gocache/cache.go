package gocache

import (
	lru2 "go-cache/gocache/lru"
	"sync"
)

type cache struct {
	lru        *lru2.Cache
	lock       sync.RWMutex
	cacheBytes uint64
}

func (c *cache) Add(key string, value ByteView) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if c.lru == nil {
		c.lru = lru2.New(c.cacheBytes, nil)
	}
	c.lru.Put(key, value)
}

func (c *cache) Get(key string) (bv ByteView, ok bool) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.lru == nil {
		return
	}

	if value, ok1 := c.lru.Get(key); ok1 {
		return value.(ByteView), true
	} else {
		return
	}

}

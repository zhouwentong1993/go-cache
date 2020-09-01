package go_cache

import (
	lru2 "go-cache/lru"
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

func (c *cache) Get(key string) (ok bool, bv ByteView) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if c.lru == nil {
		return
	}

	if ok, value := c.lru.Get(key); ok {
		return true, value.(ByteView)
	} else {
		return
	}

}

package lru

import (
	"container/list"
	"math"
)

type Cache struct {
	list      *list.List
	cache     map[string]*list.Element
	maxBytes  uint64
	currBytes uint64
	OnEvicted func(key string, value Value)
}

func New(maxBytes uint64, OnEvicted func(string, Value)) *Cache {
	if maxBytes == 0 {
		maxBytes = math.MaxUint64
	}
	return &Cache{
		list:      list.New(),
		cache:     make(map[string]*list.Element),
		maxBytes:  maxBytes,
		currBytes: 0,
		OnEvicted: OnEvicted,
	}
}

func (c Cache) Get(key string) (value Value, ok bool) {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		return kv.value, true
	} else {
		return nil, false
	}
}

func (c *Cache) Put(key string, value Value) bool {
	if ele, ok := c.cache[key]; ok {
		c.list.MoveToFront(ele)
		kv := ele.Value.(*entry)
		c.currBytes += uint64(value.Len() - kv.value.Len())
		kv.value = value
	} else {
		front := c.list.PushFront(&entry{
			key:   key,
			value: value,
		})
		c.cache[key] = front
		c.currBytes += uint64(len(key) + value.Len())
	}
	for c.currBytes > c.maxBytes {
		c.RemoveOldest()
	}
	return true
}

func (c *Cache) RemoveOldest() {
	oldestElement := c.list.Back()
	if oldestElement != nil {
		c.list.Remove(oldestElement)
		cachedData := oldestElement.Value.(*entry)
		delete(c.cache, cachedData.key)
		c.currBytes -= uint64(len(cachedData.key) + cachedData.value.Len())
		if c.OnEvicted != nil {
			c.OnEvicted(cachedData.key, cachedData.value)
		}
	}
}

type Value interface {
	Len() int
}

type entry struct {
	key   string
	value Value
}

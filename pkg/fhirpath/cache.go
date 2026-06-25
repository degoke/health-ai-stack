package fhirpath

import (
	"container/list"
	"sync"
)

type cacheEntry struct {
	key      string
	compiled CompiledExpression
}

type exprCache struct {
	mu    sync.Mutex
	size  int
	order *list.List
	items map[string]*list.Element
}

func newExprCache(size int) *exprCache {
	if size <= 0 {
		size = DefaultCacheSize
	}
	return &exprCache{
		size:  size,
		order: list.New(),
		items: make(map[string]*list.Element, size),
	}
}

func (c *exprCache) get(expr string) (CompiledExpression, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	elem, ok := c.items[expr]
	if !ok {
		return nil, false
	}
	c.order.MoveToFront(elem)
	return elem.Value.(cacheEntry).compiled, true
}

func (c *exprCache) put(expr string, compiled CompiledExpression) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if elem, ok := c.items[expr]; ok {
		elem.Value = cacheEntry{key: expr, compiled: compiled}
		c.order.MoveToFront(elem)
		return
	}
	if c.order.Len() >= c.size {
		back := c.order.Back()
		if back != nil {
			entry := back.Value.(cacheEntry)
			delete(c.items, entry.key)
			c.order.Remove(back)
		}
	}
	elem := c.order.PushFront(cacheEntry{key: expr, compiled: compiled})
	c.items[expr] = elem
}

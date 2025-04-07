package utils

import (
	"sync"
	"time"
)

type Cache struct {
	data       map[string]cacheItem
	mutex      sync.RWMutex
	defaultTTL time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		data:       make(map[string]cacheItem),
		defaultTTL: defaultTTL,
	}
}


func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.data[key]
	if !found {
		return nil, false
	}

	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]cacheItem)
}

type ParallelExecutor struct {
	wg sync.WaitGroup
}

func NewParallelExecutor() *ParallelExecutor {
	return &ParallelExecutor{}
}

func (p *ParallelExecutor) Execute(functions ...func()) {
	p.wg.Add(len(functions))

	for _, f := range functions {
		go func(fn func()) {
			defer p.wg.Done()
			fn()
		}(f)
	}

	p.wg.Wait()
}

func (p *ParallelExecutor) ExecuteWithResults(functions ...func() interface{}) []interface{} {
	results := make([]interface{}, len(functions))
	p.wg.Add(len(functions))

	for i, f := range functions {
		go func(index int, fn func() interface{}) {
			defer p.wg.Done()
			results[index] = fn()
		}(i, f)
	}

	p.wg.Wait()
	return results
}

var CommandCache = NewCache(5 * time.Minute)

package common

import (
	"sync"
	"time"
)

// Cache is a simple in-memory cache with expiration
type Cache struct {
	data       map[string]cacheItem
	mutex      sync.RWMutex
	defaultTTL time.Duration
}

type cacheItem struct {
	value      interface{}
	expiration time.Time
}

// NewCache creates a new cache with the specified default TTL
func NewCache(defaultTTL time.Duration) *Cache {
	return &Cache{
		data:       make(map[string]cacheItem),
		defaultTTL: defaultTTL,
	}
}

// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mutex.RLock()
	defer c.mutex.RUnlock()

	item, found := c.data[key]
	if !found {
		return nil, false
	}

	// Check if the item has expired
	if time.Now().After(item.expiration) {
		return nil, false
	}

	return item.value, true
}

// Set adds a value to the cache with the default TTL
func (c *Cache) Set(key string, value interface{}) {
	c.SetWithTTL(key, value, c.defaultTTL)
}

// SetWithTTL adds a value to the cache with a specific TTL
func (c *Cache) SetWithTTL(key string, value interface{}, ttl time.Duration) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data[key] = cacheItem{
		value:      value,
		expiration: time.Now().Add(ttl),
	}
}

// Delete removes a value from the cache
func (c *Cache) Delete(key string) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	delete(c.data, key)
}

// Clear removes all values from the cache
func (c *Cache) Clear() {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	c.data = make(map[string]cacheItem)
}

// ParallelExecutor executes functions in parallel and collects their results
type ParallelExecutor struct {
	wg sync.WaitGroup
}

// NewParallelExecutor creates a new parallel executor
func NewParallelExecutor() *ParallelExecutor {
	return &ParallelExecutor{}
}

// Execute executes the given functions in parallel and returns when all are complete
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

// ExecuteWithResults executes the given functions in parallel and returns their results
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

// CommandCache is a cache for command execution results
var CommandCache = NewCache(5 * time.Minute)

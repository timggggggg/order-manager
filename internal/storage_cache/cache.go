package storagecache

import (
	"fmt"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type CacheType string

const (
	DEFAULT CacheType = "default"
	LRU     CacheType = "LRU"
	LFU     CacheType = "LFU"
)

type cacheStrategy interface {
	Get(key int64) *models.Order
	Put(key int64, value *models.Order)
}

type Cache struct {
	strategy cacheStrategy
}

func NewCache(strategy cacheStrategy) *Cache {
	return &Cache{strategy: strategy}
}

func (c *Cache) SetStrategy(strat cacheStrategy) {
	c.strategy = strat
}

func (c *Cache) Get(key int64) *models.Order {
	return c.strategy.Get(key)
}

func (c *Cache) Put(key int64, value *models.Order) {
	c.strategy.Put(key, value)
}

type cacheStrategyMap map[string]func(size int64) cacheStrategy

var cacheStrategies = cacheStrategyMap{
	string(DEFAULT): func(size int64) cacheStrategy { return newDefaultCache() },
	string(LRU):     func(size int64) cacheStrategy { return newLRUCache(size) },
	string(LFU):     func(size int64) cacheStrategy { return newLFUCache(size) },
}

func NewCacheStrategy(cacheType string, cacheSize int64) (cacheStrategy, error) {
	createStrategy, exists := cacheStrategies[cacheType]
	if !exists {
		return nil, fmt.Errorf("invalid cache type: %s", cacheType)
	}
	return createStrategy(cacheSize), nil
}

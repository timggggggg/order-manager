package storagecache

import "gitlab.ozon.dev/timofey15g/homework/internal/models"

type DefaultCache struct{}

func newDefaultCache() *DefaultCache {
	return &DefaultCache{}
}

func (c *DefaultCache) Get(key int64) *models.Order {
	return nil
}

func (c *DefaultCache) Put(key int64, order *models.Order) {}

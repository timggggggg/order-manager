package postgres

import (
	"context"
	"fmt"
	"sync"
	"time"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type OrderHistoryCache struct {
	orders         models.OrdersMapStorage
	fetchFunc      func(ctx context.Context, limit, offset int64) (OrdersDBSliceStorage, error)
	updateInterval time.Duration
	LastUpdated    time.Time
	stopChan       chan struct{}
	mu             *sync.RWMutex
	Size           int64
	wg             sync.WaitGroup
}

func NewOrderHistoryCache(updateInterval time.Duration, fetchFunc func(ctx context.Context, limit, offset int64) (OrdersDBSliceStorage, error), size int64) *OrderHistoryCache {
	return &OrderHistoryCache{
		orders:         make(models.OrdersMapStorage),
		fetchFunc:      fetchFunc,
		updateInterval: updateInterval,
		stopChan:       make(chan struct{}),
		mu:             new(sync.RWMutex),
		Size:           size,
	}
}

func (c *OrderHistoryCache) StartBackgroundRefresh() {
	c.wg.Add(1)
	go c.refreshLoop()
}

func (c *OrderHistoryCache) Stop() {
	close(c.stopChan)
	c.wg.Wait()
}

func (c *OrderHistoryCache) GetHistory() models.OrdersMapStorage {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return c.orders
}

func (c *OrderHistoryCache) refreshLoop() {
	timer := time.NewTimer(c.updateInterval)

	defer c.wg.Done()
	defer timer.Stop()

	for {
		select {
		case <-timer.C:
			c.refreshCache()
			timer.Reset(c.updateInterval)
		case <-c.stopChan:
			return
		}
	}
}

func (c *OrderHistoryCache) refreshCache() {
	c.mu.Lock()
	defer c.mu.Unlock()

	orders, err := c.fetchFunc(context.Background(), c.Size, 0)
	if err != nil {
		fmt.Printf("Error updating history cache: %v", err)
		return
	}

	for _, o := range orders {
		c.orders[o.ID] = FromDTO(o)
	}

	c.LastUpdated = time.Now()
	fmt.Println("Cache updated at", c.LastUpdated)
}

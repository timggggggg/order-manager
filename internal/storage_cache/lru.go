package storagecache

import (
	"sync"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type LRUNode struct {
	key   int64
	order *models.Order
	prev  *LRUNode
	next  *LRUNode
}

type LRUCache struct {
	capacity int64
	cache    map[int64]*LRUNode
	head     *LRUNode
	tail     *LRUNode
	mu       *sync.RWMutex
}

func newLRUCache(capacity int64) *LRUCache {
	head := &LRUNode{}
	tail := &LRUNode{}
	head.next = tail
	tail.prev = head

	return &LRUCache{
		capacity: capacity,
		cache:    make(map[int64]*LRUNode),
		head:     head,
		tail:     tail,
		mu:       new(sync.RWMutex),
	}
}

func (c *LRUCache) Get(key int64) *models.Order {
	c.mu.Lock()
	defer c.mu.Unlock()

	if LRUnode, ok := c.cache[key]; ok {
		c.moveToTail(LRUnode)
		return LRUnode.order
	}
	return nil
}

func (c *LRUCache) Put(key int64, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if LRUnode, ok := c.cache[key]; ok {
		LRUnode.order = order
		c.moveToTail(LRUnode)
		return
	}

	newLRUNode := &LRUNode{
		key:   key,
		order: order,
	}

	c.cache[key] = newLRUNode
	c.addToTail(newLRUNode)

	if len(c.cache) > int(c.capacity) {
		c.removeHead()
	}
}

func (c *LRUCache) addToTail(LRUnode *LRUNode) {
	LRUnode.prev = c.tail.prev
	LRUnode.next = c.tail
	c.tail.prev.next = LRUnode
	c.tail.prev = LRUnode
}

func (c *LRUCache) removeLRUNode(LRUnode *LRUNode) {
	LRUnode.prev.next = LRUnode.next
	LRUnode.next.prev = LRUnode.prev
}

func (c *LRUCache) moveToTail(LRUnode *LRUNode) {
	c.removeLRUNode(LRUnode)
	c.addToTail(LRUnode)
}

func (c *LRUCache) removeHead() {
	if c.head.next == c.tail {
		return
	}

	oldHead := c.head.next
	c.removeLRUNode(oldHead)
	delete(c.cache, oldHead.key)
}

package storagecache

import (
	"container/list"
	"sync"

	"gitlab.ozon.dev/timofey15g/homework/internal/models"
)

type LFUNode struct {
	key   int64
	order *models.Order
	freq  int64
}

type LFUCache struct {
	capacity int64
	minFreq  int64
	keys     map[int64]*list.Element
	freqMap  map[int64]*list.List
	mu       *sync.RWMutex
}

func newLFUCache(capacity int64) *LFUCache {
	return &LFUCache{
		capacity: capacity,
		minFreq:  0,
		keys:     make(map[int64]*list.Element),
		freqMap:  make(map[int64]*list.List),
		mu:       new(sync.RWMutex),
	}
}

func (c *LFUCache) Get(key int64) *models.Order {
	c.mu.RLock()
	defer c.mu.RUnlock()

	elem, ok := c.keys[key]
	if !ok {
		return nil
	}

	LFUnode := elem.Value.(*LFUNode)
	c.removeFromFreqList(LFUnode.freq, elem)

	LFUnode.freq++
	c.addToFreqList(LFUnode, LFUnode.freq)
	c.keys[key] = elem

	return LFUnode.order
}

func (c *LFUCache) Put(key int64, order *models.Order) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if elem, ok := c.keys[key]; ok {
		LFUnode := elem.Value.(*LFUNode)
		LFUnode.order = order
		c.removeFromFreqList(LFUnode.freq, elem)
		LFUnode.freq++
		c.addToFreqList(LFUnode, LFUnode.freq)
		c.keys[key] = elem
		return
	}

	if len(c.keys) >= int(c.capacity) {
		c.evict()
	}

	newLFUNode := &LFUNode{
		key:   key,
		order: order,
		freq:  1,
	}

	c.addToFreqList(newLFUNode, 1)
	c.minFreq = 1
	c.keys[key] = c.freqMap[1].Back()
}

func (c *LFUCache) removeFromFreqList(freq int64, elem *list.Element) {
	c.freqMap[freq].Remove(elem)
	if c.freqMap[freq].Len() == 0 {
		delete(c.freqMap, freq)
		if freq == c.minFreq {
			c.minFreq++
		}
	}
}

func (c *LFUCache) addToFreqList(LFUnode *LFUNode, freq int64) {
	if _, ok := c.freqMap[freq]; !ok {
		c.freqMap[freq] = list.New()
	}
	c.freqMap[freq].PushBack(LFUnode)
}

func (c *LFUCache) evict() {
	currentList := c.freqMap[c.minFreq]
	if currentList == nil {
		return
	}

	frontElem := currentList.Front()
	if frontElem != nil {
		LFUnode := frontElem.Value.(*LFUNode)
		currentList.Remove(frontElem)
		delete(c.keys, LFUnode.key)

		if currentList.Len() == 0 {
			delete(c.freqMap, c.minFreq)
		}
	}
}

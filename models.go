package main

import (
	"errors"
	"sync"
)

// The Cache Model
type item struct {
	count byte
	next     *item
	previous *item
	expression string
	result int64
}

type order struct {
	head *item
	tail *item
}

type Cache struct {
	mx                     sync.RWMutex
	items                  map[string]*item
	size                   byte
	capacity               byte
	minCount               byte
	order
}

func NewCache(cacheSize byte) (*Cache, error) {
	if cacheSize < 2 {
		return nil, errors.New("incorrect cache size")
	}

	return &Cache{items: make(map[string]*item, cacheSize), capacity: cacheSize}, nil
}

// Increment increments the expression counter if the expression is in the cache or returns false
func (c *Cache) Increment(expression string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	_, ok = c.items[expression]
	if ok {
		c.items[expression].count++
	}
	return
}

// Add adds the new item to the cache. Trows away the oldest item unless the cache has free space
func (c *Cache) Add(expression string, result int64) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// Create item
	item := item{count: 1, expression: expression, result: result}

	// Check if we have free space
	ok = c.capacity > c.size
	if ok {
		c.size++
	} else {
		// Delete in the list the oldest item
		itemToDelete := c.order.tail
		delete(c.items, itemToDelete.expression)

		// Delete the oldest item in the order
		itemToDelete.previous.next = nil
		c.order.tail = itemToDelete.previous
	}

	// Add the new item to the cache
	c.items[expression] = &item

	// Add the new item to the order
	c.order.Add(&item)

	return
}

func (c *Cache) GetTheOldestExpression() string {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.order.head.expression
}

func (c *Cache) GetMinCount() byte {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.order.head.count
}

func (o *order) Add(item *item) {
	switch o.head {
	case nil: // The order association list is empty
		o.head = item
		o.tail = item
	default:
		item.next = o.head
		o.head.previous = item
		o.head = item
	}
}

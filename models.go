package main

import (
	"errors"
	"sync"
)

// The Cache Model

type item struct {
	Count byte
}

type expressionResult struct {
	item
	result int64
}

type orderItem struct {
	item     *item
	next     *orderItem
	previous *orderItem
}

type order struct {
	head *orderItem
	tail *orderItem
}

// The queue Model
type queue struct {
	list map[string]*item
	order
	availableSpace byte
}

type Cache struct {
	mx                     sync.RWMutex
	items                  map[string]*expressionResult
	minCountItemExpression string
	size                   byte
	capacity               byte
	minCount               byte
	queue
}

func NewCache(cacheSize, queueSize byte) (*Cache, error) {
	if cacheSize == 0 {
		return nil, errors.New("incorrect cache size")
	}

	if queueSize == 0 {
		return nil, errors.New("incorrect queue size")
	}

	newCache := Cache{items: make(map[string]*expressionResult, cacheSize), capacity: cacheSize,
		queue: queue{
			list:           make(map[string]*item, queueSize),
			order:          order{},
			availableSpace: queueSize,
		}}

	return &newCache, nil
}

// HasResult increments the expression counter if the expression is in the cache or returns false
func (c *Cache) HasResult(expression string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	_, ok = c.items[expression]
	if ok {
		c.items[expression].Count++
		if c.minCountItemExpression != "" && c.items[c.minCountItemExpression] == c.items[expression] {
			c.minCount++
		}
	}
	return
}

// CheckSpaceAndAddToCache adds an item to the cache unless the cache is full
func (c *Cache) CheckSpaceAndAddToCache(expression string, result int64, count byte) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	ok = c.capacity > c.size
	if ok {
		c.items[expression] = &expressionResult{
			result: result,
			item: item{Count: count},
		}

		if c.minCountItemExpression == "" || c.minCount > count {
			c.minCountItemExpression = expression
			c.minCount = count
		}

		c.size++
	}

	return
}

func (c *Cache) Move(expression string, result int64, count byte) (err error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// Delete item in the queue
	delete(c.queue.list, expression)

	// Pop an item in cache
	expressionToMoveToQueue := c.minCountItemExpression
	itemToMoveToQueue := &c.items[c.minCountItemExpression].item
	delete(c.items, c.minCountItemExpression)

	// Add the popped item to cache
	c.queue.list[expressionToMoveToQueue] = itemToMoveToQueue
	c.AddToQueueOrder(itemToMoveToQueue)

	// Add the input expression to the cache with the result and count
	c.items[expression] = &expressionResult{
		result: result,
		item: item{Count: count},
	}

	if c.minCount > count {
		c.minCountItemExpression = expression
		c.minCount = count
	}

	if c.minCountItemExpression == "" {
		err = errors.New("strange minCountItemExpression")
	}

	return
}

func (c *Cache) GetExpressionWithMinCount() string {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.minCountItemExpression
}

func (c *Cache) GetMinCount() byte {
	c.mx.RLock()
	defer c.mx.RUnlock()

	item, ok := c.items[c.minCountItemExpression]
	if ok {
		return item.Count
	}

	return 0
}

// HasExpression increments the expression counter if the expression is in the queue or returns false
func (c *Cache) HasExpression(expression string) byte {
	c.mx.Lock()
	defer c.mx.Unlock()

	item, ok := c.queue.list[expression]
	if ok {
		item.Count++
		return item.Count
	}

	return 0
}

func (c *Cache) AddToQueue(expression string, item *item) (err error) {
	c.mx.Lock()
	defer c.mx.Unlock()

	_, ok := c.queue.list[expression]
	if ok {
		c.queue.list[expression].Count = c.queue.list[expression].Count + 1
		return
	}

	if c.queue.availableSpace == 0 {
		delete(c.queue.list, expression)
		c.queue.list[expression] = item
		err = c.SubstituteWith(item)
	} else {
		c.queue.list[expression] = item
		c.AddToQueueOrder(item)
		c.queue.availableSpace--
	}

	return
}

func (c *Cache) AddToQueueOrder(item *item) {
	newOrderItem := &orderItem{item: item}

	switch c.queue.head {
	case nil: // list is empty
		c.queue.head = newOrderItem
		c.queue.tail = newOrderItem
	default:
		newOrderItem.next = c.queue.head
		c.queue.head.previous = newOrderItem
		c.queue.head = newOrderItem
	}
}

func (c *Cache) SubstituteWith(item *item) (err error) {
	if c.queue.head != nil {
		c.queue.head.next.previous = nil
		c.queue.tail.next = &orderItem{
			item:     item,
			next:     nil,
			previous: c.queue.tail,
		}
		c.queue.head = c.queue.head.next
	} else {
		err = errors.New("the queue is empty yet")
	}
	return
}

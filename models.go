package main

import (
	"errors"
	"sync"
)

// TMP

type ReturnItem struct {
	*item
	expression string
}

// The Cache Model

type item struct {
	Count byte
}

type expressionResult struct {
	item
	result int64
}

type Cache struct {
	mx                     sync.RWMutex
	items                  map[string]*expressionResult
	minCountItemExpression string
	size                   byte
	capacity               byte
	minCount               byte
}

func NewCache(size byte) *Cache {
	newCacheItems := make(map[string]*expressionResult, size)
	return &Cache{items: newCacheItems, capacity: size}
}

func (c *Cache) HasExpression(expression string) (ok bool) {
	c.mx.RLock()
	defer c.mx.RUnlock()

	_, ok = c.items[expression]
	if ok {
		c.items[expression].Count++
		if c.minCountItemExpression != "" && c.items[c.minCountItemExpression] == c.items[expression] {
			c.minCount++
		}
	}
	return
}

// CheckSpaceAndAdd adds an item to the cache unless the cache is full
func (c *Cache) CheckSpaceAndAdd(expression string, result int64) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	ok = c.capacity > c.size
	if ok {
		var cacheItem expressionResult
		cacheItem.result = result
		cacheItem.item.Count = 1

		c.items[expression] = &cacheItem
		c.minCountItemExpression = expression
		c.size++
	}

	return
}

func (c *Cache) Pop() (expressionToMove string, itemToMove *item) {
	expressionToMove = c.minCountItemExpression
	itemToMove = &c.items[c.minCountItemExpression].item
	delete(c.items, c.minCountItemExpression)
	c.size--
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

// The Queue Model
type Queue struct {
	mx   sync.RWMutex
	list map[string]*item
	order
	availableSpace byte
}

func (q *Queue) HasExpression(expression string) byte {
	q.mx.Lock()
	q.mx.Unlock()

	item, ok := q.list[expression]
	if ok {
		item.Count++
		return item.Count
	}

	return 0
}

func (q *Queue) Delete(expression string) {
	q.mx.Lock()
	defer q.mx.Unlock()

	delete(q.list, expression)
	q.availableSpace++
}

func (q *Queue) Move(*item) (err error) {

	return
}

func (q *Queue) Add(expression string, item *item) (err error) {
	q.mx.Lock()
	defer q.mx.Unlock()

	_, ok := q.list[expression]
	if ok {
		q.list[expression].Count = q.list[expression].Count + 1
		return
	}

	if q.availableSpace == 0 {
		delete(q.list, expression)
		q.list[expression] = item
		err = q.order.SubstituteWith(item)
	} else {
		q.list[expression] = item
		q.order.Add(item)
		q.availableSpace--
	}
	return
}

func NewQueue(size byte) *Queue {
	newQueueList := make(map[string]*item, size)
	return &Queue{list: newQueueList, availableSpace: size}
}

type orderItem struct {
	item     *item
	next     *orderItem
	previous *orderItem
}

type order struct {
	mx   sync.RWMutex
	head *orderItem
	tail *orderItem
}

func (o *order) Add(item *item) {
	o.mx.Lock()
	defer o.mx.Unlock()

	newOrderItem := &orderItem{item: item}

	switch o.head {
	case nil: // list is empty
		o.head = newOrderItem
		o.tail = newOrderItem
	default:
		newOrderItem.next = o.head
		o.head.previous = newOrderItem
		o.head = newOrderItem
	}
}

func (o *order) SubstituteWith(item *item) (err error) {
	o.mx.Lock()
	defer o.mx.Unlock()

	if o.head != nil {
		o.head.next.previous = nil
		o.tail.next = &orderItem{
			item:     item,
			next:     nil,
			previous: o.tail,
		}
		o.head = o.head.next
	} else {
		err = errors.New("the queue is empty yet")
	}
	return
}

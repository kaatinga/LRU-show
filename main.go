package LRU

import (
	"errors"
	"sync"
)

// The LRU Cache Item Model
type item struct {
	count    byte
	next     *item
	previous *item
	index    string
	data     interface{}
}

// The LRU Cache Order SubModel
type order struct {
	head *item
	tail *item
}

// The LRU Cache Model
type Cache struct {
	mx       sync.RWMutex
	items    map[string]*item
	size     byte
	capacity byte
	order
}

func NewCache(cacheSize byte) (*Cache, error) {
	if cacheSize < 2 {
		return nil, errors.New("incorrect cache size")
	}

	return &Cache{items: make(map[string]*item, cacheSize), capacity: cacheSize}, nil
}

// Increment increments the expression counter if an item with such an index exists in the cache or returns false
func (c *Cache) Increment(index string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	var gottenItem *item
	gottenItem, ok = c.items[index]
	if ok {
		gottenItem.count++

		if c.order.head != gottenItem {

			// Set prev. and next fields for the items around
			if c.order.tail != gottenItem {
				gottenItem.previous.next, gottenItem.next.previous = gottenItem.next, gottenItem.previous
			} else {
				gottenItem.previous.next = nil
			}

			// Move the item to the beginning of the order
			gottenItem.previous = nil
			gottenItem.next = c.order.head

			c.order.head.previous = gottenItem
			c.order.head = gottenItem
		}
	}
	return
}

// Delete deletes an item Cache with the index in the signature
func (c *Cache) Delete(index string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	_, ok = c.items[index]
	if ok {
		c.items[index].previous = nil
		c.items[index].next = nil
	} else {
		return
	}

	// In case it was the only item in the Cache
	if c.items[index] == c.order.head && c.items[index] == c.order.tail {
		c.order.head = nil
		c.order.tail = nil
	}

	if c.items[index].previous != c.order.head {
		if c.items[index].next != c.order.tail {
			c.items[index].previous.next = c.items[index].next
		} else {
			c.items[index].previous.next = nil
		}
	}

	delete(c.items, index)
	return
}

// Add adds the new item to the Cache. Trows away the oldest item unless the Cache has free space
func (c *Cache) Add(index string, data interface{}) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	// New item creation
	item := item{count: 1, index: index, data: data}

	// Check if we have free space
	ok = c.capacity > c.size
	if ok {
		c.size++
	} else {
		// Delete in the list the oldest item
		itemToDelete := c.order.tail
		delete(c.items, itemToDelete.index)

		// Delete the oldest item in the order
		itemToDelete.previous.next = nil
		c.order.tail = itemToDelete.previous
	}

	// add the new item to the cache
	c.items[index] = &item

	// add the new item to the order
	c.order.add(&item)

	return
}

// GetTheOldestIndex returns the oldest index in the cache
func (c *Cache) GetTheOldestIndex() string {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.order.tail.index
}

// GetMinCount returns the oldest index count field value
func (c *Cache) GetMinCount() byte {
	c.mx.RLock()
	defer c.mx.RUnlock()

	return c.order.tail.count
}

// add is an internal package method to keep order of the items
func (o *order) add(item *item) {
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

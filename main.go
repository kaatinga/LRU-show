package LRU

import (
<<<<<<< HEAD
	"errors"
	"sync"
=======
	"github.com/gdamore/tcell"
	"github.com/kaatinga/calc"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"

	"github.com/kaatinga/LRU"
>>>>>>> updated to use https://github.com/kaatinga/LRU
)

// The LRU Cache Item Model
type item struct {
	count    byte
	next     *item
	previous *item
	index    string
	data     interface{}
}

<<<<<<< HEAD
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
=======
const (
	capacity byte = 5
)

func main() {

	// Create a new cache and queue
	var (
		Show LRUMonitor
		err  error
	)

	Show.cache, err = LRU.NewCache(capacity)
	if err != nil {
		log.Fatalln(err)
>>>>>>> updated to use https://github.com/kaatinga/LRU
	}

	return &Cache{items: make(map[string]*item, cacheSize), capacity: cacheSize}, nil
}

// Increment increments the expression counter if an item with such an index exists in the cache or returns false
func (c *Cache) Increment(index string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

<<<<<<< HEAD
	var gottenItem *item
	gottenItem, ok = c.items[index]
	if ok {
		gottenItem.count++
=======
	Show.cacheData = tview.NewTable()
	Show.orderData = tview.NewTable()
>>>>>>> updated to use https://github.com/kaatinga/LRU

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

<<<<<<< HEAD
			c.order.head.previous = gottenItem
			c.order.head = gottenItem
		}
	}
	return
}
=======
	grid.AddItem(Show.cacheData, 1, 0, 1, 1, 0, 0, false).
		AddItem(Show.orderData, 1, 1, 1, 1, 0, 0, false).
		AddItem(messageLog, 1, 2, 1, 1, 0, 0, false).
		AddItem(inputField, 2, 0, 1, 3, 0, 0, true)
>>>>>>> updated to use https://github.com/kaatinga/LRU

// Delete deletes an item Cache with the index in the signature
func (c *Cache) Delete(index string) (ok bool) {
	c.mx.Lock()
	defer c.mx.Unlock()

	var gottenItem *item
	gottenItem, ok = c.items[index]
	if !ok {
		return
	}

	if gottenItem.previous != nil {
		if gottenItem.next != nil {
			gottenItem.previous.next, gottenItem.next.previous = gottenItem.next, gottenItem.previous
		} else {
			gottenItem.previous.next = nil
			c.order.tail = gottenItem.previous
		}
	}

	if gottenItem.next != nil {
		if gottenItem.previous == nil {
			gottenItem.next.previous = nil
			c.order.head = gottenItem.next
		}
	}

	gottenItem.previous = nil
	gottenItem.next = nil

<<<<<<< HEAD
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
=======
			// AddToQueue a message to the log
			messageRow = AddMessage(messageLog, strings.Join([]string{"The user entered expression", expression}, ": "), messageRow)

			// Start to work with cache
			if !Show.cache.Increment(expression) {

				// Calculate the result
				var result int64
				result, err = calc.Calc(expression)
				if err != nil {
					messageRow = AddMessage(messageLog, "Error: "+err.Error(), messageRow)
					break
				}
				messageRow = AddMessage(messageLog, strings.Join([]string{"The expression result was calculated", strconv.Itoa(int(result))}, ": "), messageRow)

				if !Show.cache.Add(expression, result) {
					messageRow = AddMessage(messageLog, "The cache had no free space", messageRow)
					messageRow = AddMessage(messageLog, "The oldest item was substituted with the new one", messageRow)
				} else {
					messageRow = AddMessage(messageLog, "The expression was added to the cache. The cache was not full", messageRow)
				}
			} else {
				messageRow = AddMessage(messageLog, "The result was found in the cache", messageRow)
			}

			// print the cache onscreen
			printMessage := Show.PrintCache()
			if printMessage != "" {
				messageRow = AddMessage(messageLog, printMessage, messageRow)
			}
>>>>>>> updated to use https://github.com/kaatinga/LRU

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

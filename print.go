package main

import (
	"github.com/rivo/tview"
	"strconv"
)

type LRUMonitor struct {
	cacheData *tview.Table
	queueData *tview.Table
	orderData *tview.Table
	cache     *Cache
}

func (LRU *LRUMonitor) PrintLRU() (message string) {

	var i int

	LRU.cacheData.Clear()
	LRU.cacheData.SetCellSimple(i, 0, "Expression")
	LRU.cacheData.SetCellSimple(i, 1, "Result")
	LRU.cacheData.SetCellSimple(i, 2, "Count")
	i++
	for expression, cacheItem := range LRU.cache.items {
		LRU.cacheData.SetCellSimple(i, 0, expression)
		LRU.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(cacheItem.result)))
		LRU.cacheData.SetCellSimple(i, 2, strconv.Itoa(int(cacheItem.count)))
		i++
	}
	LRU.cacheData.SetCellSimple(i, 0, "----------")
	i++
	LRU.cacheData.SetCellSimple(i, 0, "Oldest: ")
	i++
	LRU.cacheData.SetCellSimple(i, 1, "Expr.: ")
	LRU.cacheData.SetCellSimple(i, 2, LRU.cache.GetTheOldestExpression())
	i++
	LRU.cacheData.SetCellSimple(i, 1, "Count: ")
	LRU.cacheData.SetCellSimple(i, 2, strconv.Itoa(int(LRU.cache.GetMinCount())))
	i++
	LRU.cacheData.SetCellSimple(i, 0, "Capacity: ")
	LRU.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(LRU.cache.capacity)))
	i++
	LRU.cacheData.SetCellSimple(i, 0, "Size: ")
	LRU.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(LRU.cache.size)))

	// Reuse i
	i = 0
	LRU.orderData.Clear()
	LRU.orderData.SetCellSimple(i, 0, "Expression")
	LRU.orderData.SetCellSimple(i, 1, "Result")
	LRU.orderData.SetCellSimple(i, 2, "Count")
	LRU.orderData.SetCellSimple(i, 3, "Prev.")
	LRU.orderData.SetCellSimple(i, 4, "Next")
	i++
	var orderItem *item
	orderItem = LRU.cache.order.head
	if orderItem != nil {
		for {
			LRU.orderData.SetCellSimple(i, 0, orderItem.expression)
			LRU.orderData.SetCellSimple(i, 1, strconv.Itoa(int(orderItem.result)))
			LRU.orderData.SetCellSimple(i, 2, strconv.Itoa(int(orderItem.count)))

			if orderItem.previous != nil {
				LRU.orderData.SetCellSimple(i, 3, orderItem.previous.expression)
			}

			if orderItem.next != nil {
				LRU.orderData.SetCellSimple(i, 4, orderItem.next.expression)
			}

			i++
			if orderItem.next == nil {
				return
			}

			orderItem = orderItem.next
		}
	}
	return "must not happen"
}

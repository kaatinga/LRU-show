package main

import (
	"github.com/rivo/tview"
	"strconv"
)

type LRUMonitor struct {
	cacheData *tview.Table
	queueData *tview.Table
	queueOrderData *tview.Table
	cache *Cache
}

func (LRU *LRUMonitor) PrintLRU() {

	var i int
	LRU.cacheData.Clear()
	LRU.cacheData.SetCellSimple(i, 0, "Expression")
	LRU.cacheData.SetCellSimple(i, 1, "Result")
	LRU.cacheData.SetCellSimple(i, 2, "Count")
	i++
	for expression, cacheItem := range LRU.cache.items {
		LRU.cacheData.SetCellSimple(i, 0, expression)
		LRU.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(cacheItem.result)))
		LRU.cacheData.SetCellSimple(i, 2, strconv.Itoa(int(cacheItem.Count)))
		i++
	}
	LRU.cacheData.SetCellSimple(i, 0, "----------")
	i++
	LRU.cacheData.SetCellSimple(i, 0, "Min Expr.: ")
	LRU.cacheData.SetCellSimple(i, 1, LRU.cache.GetExpressionWithMinCount())
	i++
	LRU.cacheData.SetCellSimple(i, 0, "Min Count: ")
	LRU.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(LRU.cache.GetMinCount())))
	i++
	LRU.cacheData.SetCellSimple(i, 0, "Capacity: ")
	LRU.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(LRU.cache.capacity)))
	i++
	LRU.cacheData.SetCellSimple(i, 0, "Size: ")
	LRU.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(LRU.cache.size)))

	// print the queue onscreen
	LRU.queueData.Clear()
	i = 0
	LRU.queueData.SetCellSimple(i, 0, "Expression")
	LRU.queueData.SetCellSimple(i, 0, "Count")
	i++
	for expression, queueItem := range LRU.cache.queue.list {
		LRU.queueData.SetCellSimple(i, 0, expression)
		LRU.queueData.SetCellSimple(i, 1, strconv.Itoa(int(queueItem.Count)))
		i++
	}
	LRU.queueData.SetCellSimple(i, 0, "----------")
	i++
	LRU.queueData.SetCellSimple(i, 0, "Av. space: ")
	LRU.queueData.SetCellSimple(i, 1, strconv.Itoa(int(LRU.cache.queue.availableSpace)))
}

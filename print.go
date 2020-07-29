package main

import (
	"github.com/kaatinga/LRU"
	"github.com/rivo/tview"
	"strconv"
)

type LRUMonitor struct {
	cacheData *tview.Table
	queueData *tview.Table
	orderData *tview.Table
	cache     *LRU.Cache
}

func (display *LRUMonitor) PrintCache() (message string) {

	var i int

	display.cacheData.SetCellSimple(i, 0, "Oldest: ")
	i++
	display.cacheData.SetCellSimple(i, 1, "Expr.: ")
	display.cacheData.SetCellSimple(i, 2, display.cache.GetTheOldestIndex())
	i++
	display.cacheData.SetCellSimple(i, 1, "Count: ")
	display.cacheData.SetCellSimple(i, 2, strconv.Itoa(int(display.cache.GetTheOldestCount())))
	i++
	display.cacheData.SetCellSimple(i, 0, "Capacity: ")
	display.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(capacity)))
	i++
	display.cacheData.SetCellSimple(i, 0, "Size: ")
	display.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(display.cache.GetTheCacheSize())))
	i++
	display.cacheData.SetCellSimple(i, 0, "Head: ")
	display.cacheData.SetCellSimple(i, 1, display.cache.GetTheHeadIndex())
	i++
	display.cacheData.SetCellSimple(i, 0, "Tail: ")
	display.cacheData.SetCellSimple(i, 1, display.cache.GetTheOldestIndex())
	// Reuse i
	i = 0
	display.orderData.Clear()
	display.orderData.SetCellSimple(i, 0, "Index")
	display.orderData.SetCellSimple(i, 1, "Result")
	display.orderData.SetCellSimple(i, 2, "Count")
	display.orderData.SetCellSimple(i, 3, "Prev.")
	display.orderData.SetCellSimple(i, 4, "Next")
	i++
	var currentIndex, nextIndex, previousIndex string
	currentIndex = display.cache.GetTheHeadIndex()
	var result interface{}
	if currentIndex != "" {
		for {

			result, _ = display.cache.GetStoredData(currentIndex)

			display.orderData.SetCellSimple(i, 0, currentIndex)

			if result != nil {
				display.orderData.SetCellSimple(i, 1, strconv.Itoa(int(result.(int64))))
			}

			display.orderData.SetCellSimple(i, 2, "-")

			previousIndex = display.cache.GetThePreviousItemIndex(currentIndex)
			if previousIndex != "" {
				display.orderData.SetCellSimple(i, 3, previousIndex)
			}

			nextIndex = display.cache.GetTheNextItemIndex(currentIndex)
			if nextIndex != "" {
				display.orderData.SetCellSimple(i, 4, nextIndex)
			}

			i++
			if nextIndex == "" {
				return
			}

			currentIndex = nextIndex
		}
	}
	return "must not happen"
}

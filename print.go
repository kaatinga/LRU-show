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
	display.cacheData.SetCellSimple(i, 2, strconv.Itoa(int(display.cache.GetMinCount())))
	i++
	display.cacheData.SetCellSimple(i, 0, "Capacity: ")
	display.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(capacity)))
	i++
	display.cacheData.SetCellSimple(i, 0, "Size: ")
	display.cacheData.SetCellSimple(i, 1, strconv.Itoa(int(display.cache.GetTheCacheSize())))

	// Reuse i
	i = 0
	display.orderData.Clear()
	display.orderData.SetCellSimple(i, 0, "Index")
	display.orderData.SetCellSimple(i, 1, "Result")
	display.orderData.SetCellSimple(i, 2, "Count")
	display.orderData.SetCellSimple(i, 3, "Prev.")
	display.orderData.SetCellSimple(i, 4, "Next")
	i++
	var currentIndex, nextIndex string
	currentIndex = display.cache.GetTheHeadIndex()
	if currentIndex != "" {
		for {
			display.orderData.SetCellSimple(i, 0, currentIndex)
			display.orderData.SetCellSimple(i, 1, "-")
			display.orderData.SetCellSimple(i, 2, "-")

			//if currentIndex.previous != nil {
			//	display.orderData.SetCellSimple(i, 3, currentIndex.previous.expression)
			//}

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

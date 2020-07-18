package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"log"
	"strings"
	"time"
)

type CacheItem struct {
	Count     byte
	LastUsage time.Time
}

type Cache struct {
	Items [4]CacheItem
}

func NewCache() *Cache {
	return &Cache{Items: [4]CacheItem{}}
}

type Queue struct {
	List map[string]CacheItem
}

func NewQueue() *Queue {
	return &Queue{List: make(map[string]CacheItem, 100)}
}

func main() {

	// Create a new cache
	var cache = NewCache()
	var queue = NewQueue()

	// Announce the app
	app := tview.NewApplication()

	newPrimitive := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	cacheData := newPrimitive("Cache")
	queueData := newPrimitive("Queue")
	messageLog := tview.NewTable()
	inputField := tview.NewInputField().
		SetLabel("Enter a math expression (press ESC to exit): ").
		SetPlaceholder("1 + 2").
		SetFieldWidth(0)

	// Grid Layout
	grid := tview.NewGrid().
		SetRows(0, 1).
		SetColumns(30, 30, 0).
		SetBorders(true)

	grid.AddItem(cacheData, 0, 0, 1, 1, 0, 0, false).
		AddItem(queueData, 0, 1, 1, 1, 0, 0, false).
		AddItem(messageLog, 0, 2, 1, 1, 0, 0, false).
		AddItem(inputField, 1, 0, 1, 3, 0, 0, true)

	// submitted is toggled each time Enter is pressed
	//var submitted bool

	var messageRow int

	// Capture user input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyEnter:

			// submitted = !submitted
			// if submitted {
			expression := strings.Trim(inputField.GetText(), " ")
			if strings.TrimSpace(expression) == "" {
				return event
			}

			var queueItem CacheItem
			var ok bool
			queueItem, ok = queue.List[expression]
			if ok {

				queue.List[expression].Count = queue.List[expression].Count + 1

				for key, cacheItem := range cache.Items {
					if cacheItem.Count < queueItem.Count {
						cache.Items[key] = queueItem
					}
				}

			} else {
				queue.List[expression] = CacheItem{}

				for key, cacheItem := range cache.Items {
					if cacheItem.Count < queueItem.Count {
						if cache.Items[key].Count == 0 {
							cache.Items[key] = cacheItem
						}
					}
				}

			}

			// add a message to the log
			messageLog.SetCellSimple(messageRow, 0, expression)
			messageRow++
			messageLog.ScrollToEnd()

			//	// Create a modal dialog
			//	m := tview.NewModal().
			//		SetText(fmt.Sprintf("You entered, %s!", expression)).
			//		AddButtons([]string{"Ok"})
			//
			//	// Display and focus the dialog
			//	app.SetRoot(m, true).SetFocus(m)
			//} else {
			// Clear the input field

			inputField.SetText("")

			// Display appGrid and focus the input field
			app.SetRoot(grid, true).SetFocus(inputField)

			//}

			return nil
		case tcell.KeyEsc:
			// Exit the application
			app.Stop()
			return nil
		}

		return event
	})

	err := app.SetRoot(grid, true).Run()
	if err != nil {
		log.Println(err)
	}
}

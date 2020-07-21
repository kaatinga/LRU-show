package main

import (
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/kaatinga/calc"
)

const (
	FormattedDateTime = "02.01.2006 15:04"
)

var (
	moscow *time.Location
)

func AddMessage(messageLog *tview.Table, text string, messageRow int) int {
	messageLog.SetCellSimple(messageRow, 0, text)
	messageLog.ScrollToEnd()
	messageRow++
	return messageRow
}

func main() {

	// set location
	moscow, _ = time.LoadLocation("Europe/Moscow")

	// Create a new cache
	var queue = NewQueue(100)
	var cache = NewCache(4)

	// Announce the app
	app := tview.NewApplication()

	title := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	cacheData := tview.NewTable()
	queueData := tview.NewTable()
	queueOrderData := tview.NewTable()
	messageLog := tview.NewTable()
	inputField := tview.NewInputField().
		SetLabel("Enter a math expression (press ESC to exit): ").
		SetPlaceholder("1 + 2").
		SetFieldWidth(0)

	// Grid Layout
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(35, 25, 25, 0).
		SetBorders(true)

	// titles
	grid.AddItem(title("=== Cache ==="), 0, 0, 1, 1, 0, 0, false).
		AddItem(title("=== Queue ==="), 0, 1, 1, 1, 0, 0, false).
		AddItem(title("=== Queue Order ==="), 0, 2, 1, 1, 0, 0, false).
		AddItem(title("=== Message Log ==="), 0, 3, 1, 1, 0, 0, false)

	grid.AddItem(cacheData, 1, 0, 1, 1, 0, 0, false).
		AddItem(queueData, 1, 1, 1, 1, 0, 0, false).
		AddItem(queueOrderData, 1, 2, 1, 1, 0, 0, false).
		AddItem(messageLog, 1, 3, 1, 1, 0, 0, false).
		AddItem(inputField, 2, 0, 1, 4, 0, 0, true)

	// submitted is toggled each time Enter is pressed
	// var submitted bool

	var messageRow int

	// Capture user input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {

		// Anything handled here will be executed on the main thread
		switch event.Key() {
		case tcell.KeyEnter:

			// submitted = !submitted
			// if submitted {
			expression := strings.ReplaceAll(inputField.GetText(), " ", "")
			if expression == "" {
				return event
			}

			// Add a message to the log
			messageRow = AddMessage(messageLog, strings.Join([]string{"The user entered expression", expression}, ": "), messageRow)

			// Start to work with cache and queue

			if !cache.HasExpression(expression) {

				// Calculate the result
				result, err := calc.Calc(expression)
				if err != nil {
					messageRow = AddMessage(messageLog, err.Error(), messageRow)
					break
				}
				messageRow = AddMessage(messageLog, strings.Join([]string{"The expression result was calculated", strconv.Itoa(int(result))}, ": "), messageRow)

				if !cache.CheckSpaceAndAdd(expression, result) {

					// Work with queue
					queueHas := queue.HasExpression(expression)
					if queueHas != 0 {
						if queueHas > cache.GetMinCount() {

							// Move the least popular item to queue
							queue.Delete(expression)
							expressionToMove, itemToMove := cache.Pop()
							err = queue.Add(expressionToMove, itemToMove)
							if err != nil {
								messageRow = AddMessage(messageLog, err.Error(), messageRow)
								break
							}

							// Add to the cache the new item
							ok := cache.CheckSpaceAndAdd(expression, result)
							if ok {
								messageRow = AddMessage(messageLog, strings.Join([]string{"The expression was moved to the cache", expression}, ": "), messageRow)
							}
						}
					} else { // The Queue has no such an expression. The expression is new!
						err = queue.Add(expression, &item{Count: 1})
						if err != nil {
							messageRow = AddMessage(messageLog, err.Error(), messageRow)
							break
						}
					}
				}
			} else {
				messageRow = AddMessage(messageLog, strings.Join([]string{"The result was found in the cache", strconv.Itoa(int(cache.items[expression].result))}, ": "), messageRow)
			}

			// print the cache onscreen
			var i int
			for expression, cacheItem := range cache.items {
				cacheData.SetCellSimple(i, 0, strings.Join([]string{"Expr.", expression}, ": "))
				cacheData.SetCellSimple(i, 1, strings.Join([]string{"result", strconv.Itoa(int(cacheItem.result))}, ": "))
				cacheData.SetCellSimple(i, 2, strings.Join([]string{"Count", strconv.Itoa(int(cacheItem.Count))}, ": "))
				i++
			}
			cacheData.SetCellSimple(i, 0, "---")
			cacheData.SetCellSimple(i, 1, "   ")
			cacheData.SetCellSimple(i, 2, "   ")
			i++
			cacheData.SetCellSimple(i, 0, "Min Expr.: ")
			cacheData.SetCellSimple(i, 2, cache.GetExpressionWithMinCount())
			i++
			cacheData.SetCellSimple(i, 0, "Min Count: ")
			cacheData.SetCellSimple(i, 2, strconv.Itoa(int(cache.GetMinCount())))
			i++
			cacheData.SetCellSimple(i, 0, "Capacity: ")
			cacheData.SetCellSimple(i, 2, strconv.Itoa(int(cache.capacity)))
			i++
			cacheData.SetCellSimple(i, 0, "Min Expr.: ")
			cacheData.SetCellSimple(i, 2, cache.minCountItemExpression)
			i++
			cacheData.SetCellSimple(i, 0, "Size: ")
			cacheData.SetCellSimple(i, 2, strconv.Itoa(int(cache.size)))

			// print the queue onscreen
			i = 0
			for expression, cacheItem := range queue.list {
				queueData.SetCellSimple(i, 0, strings.Join([]string{"Expr.", expression}, ": "))
				queueData.SetCellSimple(i, 1, strings.Join([]string{"Count", strconv.Itoa(int(cacheItem.Count))}, ": "))
				i++
			}
			queueData.SetCellSimple(i, 0, "---")
			queueData.SetCellSimple(i, 1, "   ")
			queueData.SetCellSimple(i, 2, "   ")
			i++
			queueData.SetCellSimple(i, 0, "Av. space: ")
			queueData.SetCellSimple(i, 1, strconv.Itoa(int(queue.availableSpace)))

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

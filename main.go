package main

import (
	"github.com/gdamore/tcell"
	"github.com/kaatinga/calc"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
)

func AddMessage(messageLog *tview.Table, text string, messageRow int) int {
	messageLog.SetCellSimple(messageRow, 0, text)
	messageLog.ScrollToEnd()
	messageRow++
	return messageRow
}

func main() {

	// Create a new cache and queue
	var (
		LRU LRUMonitor
		err error
	)

	LRU.cache, err = NewCache(5)
	if err != nil {
		log.Fatalln(err)
	}

	// Announce the app
	app := tview.NewApplication()

	title := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	LRU.cacheData = tview.NewTable()
	LRU.orderData = tview.NewTable()

	messageLog := tview.NewTable()
	inputField := tview.NewInputField().
		SetLabel("Enter a math expression (press ESC to exit): ").
		SetPlaceholder("1 + 2").
		SetFieldWidth(0)

	// Grid Layout
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(25, 40, 0).
		SetBorders(true)

	// titles
	grid.AddItem(title("== Cache (No order) =="), 0, 0, 1, 1, 0, 0, false).
		AddItem(title("== Order (The oldest in the down) =="), 0, 1, 1, 1, 0, 0, false).
		AddItem(title("== Message Log =="), 0, 2, 1, 1, 0, 0, false)

	grid.AddItem(LRU.cacheData, 1, 0, 1, 1, 0, 0, false).
		AddItem(LRU.orderData, 1, 1, 1, 1, 0, 0, false).
		AddItem(messageLog, 1, 2, 1, 1, 0, 0, false).
		AddItem(inputField, 2, 0, 1, 3, 0, 0, true)

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

			// AddToQueue a message to the log
			messageRow = AddMessage(messageLog, strings.Join([]string{"The user entered expression", expression}, ": "), messageRow)

			// Start to work with cache
			if !LRU.cache.Increment(expression) {

				// Calculate the result
				var result int64
				result, err = calc.Calc(expression)
				if err != nil {
					messageRow = AddMessage(messageLog, "Error: "+err.Error(), messageRow)
					break
				}
				messageRow = AddMessage(messageLog, strings.Join([]string{"The expression result was calculated", strconv.Itoa(int(result))}, ": "), messageRow)

				if !LRU.cache.Add(expression, result) {
					messageRow = AddMessage(messageLog, "The cache had no free space", messageRow)
					messageRow = AddMessage(messageLog, "The oldest item was substituted with the new one", messageRow)
				} else {
					messageRow = AddMessage(messageLog, "The expression was added to the cache. The cache was not full", messageRow)
				}
			} else {
				messageRow = AddMessage(messageLog, strings.Join([]string{"The result was found in the cache", strconv.Itoa(int(LRU.cache.items[expression].result))}, ": "), messageRow)
			}

			// print the cache onscreen
			printMessage := LRU.PrintLRU()
			if printMessage != "" {
				messageRow = AddMessage(messageLog, printMessage, messageRow)
			}

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

	err = app.SetRoot(grid, true).Run()
	if err != nil {
		log.Println(err)
	}
}

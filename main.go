package main

import (
	"errors"
	"github.com/gdamore/tcell"
	"github.com/rivo/tview"
	"log"
	"strconv"
	"strings"
	"time"
)

const (
	FormattedDateTime = "02.01.2006 15:04"
)

var (
	moscow *time.Location
)

type CacheItem struct {
	Count     byte
	LastUsage time.Time
}

type ExpressionResult struct {
	CacheItem
	Result int64
}

type Cache struct {
	Items map[string]ExpressionResult
}

func NewCache() *Cache {
	return &Cache{Items: make(map[string]ExpressionResult, 4)}
}

type Queue struct {
	List map[string]CacheItem
}

func NewQueue() *Queue {
	return &Queue{List: make(map[string]CacheItem, 100)}
}

func AddMessage(messageLog *tview.Table, text string, messageRow int) int {
	messageLog.SetCellSimple(messageRow, 0, text)
	messageLog.ScrollToEnd()
	messageRow++
	return messageRow
}

func findNums(e string) (nums map[string]int64, sign string, err error) {
	nums = make(map[string]int64, 2)
	currentStep := "num1"
	var num int64
	for _, elem := range e {
		switch elem {
		case '*':
			sign = "*"
			nums[currentStep] = num
			currentStep = "num2"
			num = 0
		case '/':
			sign = "/"
			nums[currentStep] = num
			currentStep = "num2"
			num = 0
		case '+':
			sign = "+"
			nums[currentStep] = num
			currentStep = "num2"
			num = 0
		case '-':
			sign = "-"
			nums[currentStep] = num
			currentStep = "num2"
			num = 0
		default:
			if elem < 48 || elem > 57 {
				err = errors.New("incorrect expression 2")
				return
			}
			num = num*10 + int64(elem) - 48
		}

		// set num2
		nums[currentStep] = num
	}
	return
}

func Calc(e string) (result int64, err error) {
	var nums map[string]int64
	var sign string
	nums, sign, err = findNums(e)
	if err != nil {
		return
	}

	switch sign {
	case "*":
		result = nums["num1"] * nums["num2"]
	case "/":
		if nums["num2"] == 0 {
			return 0, errors.New("you tried to divide by zero")
		}
		result = nums["num1"] / nums["num2"]
	case "+":
		result = nums["num1"] + nums["num2"]
	case "-":
		result = nums["num1"] - nums["num2"]
	default:
		return 0.0, errors.New("incorrect expression 1")
	}

	return
}

func main() {

	// set location
	moscow, _ = time.LoadLocation("Europe/Moscow")

	// Create a new cache
	var cache = NewCache()
	var queue = NewQueue()

	// Announce the app
	app := tview.NewApplication()

	title := func(text string) tview.Primitive {
		return tview.NewTextView().
			SetTextAlign(tview.AlignCenter).
			SetText(text)
	}

	cacheData := tview.NewTable()
	queueData := tview.NewTable()
	messageLog := tview.NewTable()
	inputField := tview.NewInputField().
		SetLabel("Enter a math expression (press ESC to exit): ").
		SetPlaceholder("1 + 2").
		SetFieldWidth(0)

	// Grid Layout
	grid := tview.NewGrid().
		SetRows(1, 0, 1).
		SetColumns(60, 50, 0).
		SetBorders(true)

	// titles
	grid.AddItem(title("=== Cache ==="), 0, 0, 1, 1, 0, 0, false).
		AddItem(title("=== Queue ==="), 0, 1, 1, 1, 0, 0, false).
		AddItem(title("=== Message Log ==="), 0, 2, 1, 1, 0, 0, false)

	grid.AddItem(cacheData, 1, 0, 1, 1, 0, 0, false).
		AddItem(queueData, 1, 1, 1, 1, 0, 0, false).
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
			expression := strings.Trim(inputField.GetText(), " ")
			if expression == "" {
				return event
			}

			// Add a message to the log
			messageRow = AddMessage(messageLog, strings.Join([]string{"The user entered expression", expression}, ": "), messageRow)

			// Start to work with cache and queue
			var expressionInCache ExpressionResult
			var ok bool

			// Look for the Result in the cache
			expressionInCache, ok = cache.Items[expression]
			if ok {
				messageRow = AddMessage(messageLog, strings.Join([]string{"The expression result was found in the cache", strconv.Itoa(int(expressionInCache.Result))}, ": "), messageRow)
				expressionInCache.Count = expressionInCache.Count + 1
				cache.Items[expression] = expressionInCache
			} else {

				// Calculate the result
				result, err := Calc(expression)
				if err != nil {
					messageRow = AddMessage(messageLog, err.Error(), messageRow)
					break
				}

				messageRow = AddMessage(messageLog, strings.Join([]string{"The expression result was calculated", strconv.Itoa(int(result))}, ": "), messageRow)

				// Look for the expression in the queue
				var currentItem CacheItem // temporary CacheItem
				currentItem, ok = queue.List[expression]
				if ok {

					// increase the weight of the expression
					currentItem.Count = currentItem.Count + 1
					currentItem.LastUsage = time.Now()

					// return cacheItem to the queue map
					queue.List[expression] = currentItem

				} else {

					// TODO: only if map has free capacity
					currentItem.Count = 1
					currentItem.LastUsage = time.Now()
					queue.List[expression] = currentItem
				}

				if len(cache.Items) < 4 {
					cache.Items[expression] = ExpressionResult{
						CacheItem: currentItem,
						Result:    result,
					}
					delete(queue.List, expression) // Delete the cacheItem from the queue
				} else {

					// TODO: use binary search to substitute a CacheItem with currentItem
					// Update cache if there is a cacheItem that is less popular
					for key, cacheItem := range cache.Items {
						if cacheItem.Count < currentItem.Count {
							delete(cache.Items, key)       // Delete a less popular cacheItem
							delete(queue.List, expression) // Delete the cacheItem from the queue
							cache.Items[expression] = ExpressionResult{
								CacheItem: currentItem,
								Result:    result,
							}
						}
					}
				}
			}

			// print the cache onscreen
			var i int
			for expression, cacheItem := range cache.Items {
				cacheData.SetCellSimple(i, 0, strings.Join([]string{"Expr.", expression}, ": "))
				cacheData.SetCellSimple(i, 1, strings.Join([]string{"Result", strconv.Itoa(int(cacheItem.Result))}, ": "))
				cacheData.SetCellSimple(i, 2, strings.Join([]string{"Weight", strconv.Itoa(int(cacheItem.Count))}, ": "))
				cacheData.SetCellSimple(i, 3, strings.Join([]string{"Last Usage",cacheItem.LastUsage.In(moscow).Format(FormattedDateTime)}, ": "))
				i++
			}

			// print the queue onscreen
			i = 0
			for expression, cacheItem := range queue.List {
				queueData.SetCellSimple(i, 0, strings.Join([]string{"Expr.", expression}, ": "))
				queueData.SetCellSimple(i, 1, strings.Join([]string{"Weight", strconv.Itoa(int(cacheItem.Count))}, ": "))
				queueData.SetCellSimple(i, 2, strings.Join([]string{"Last Usage",cacheItem.LastUsage.In(moscow).Format(FormattedDateTime)}, ": "))
				i++
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

	err := app.SetRoot(grid, true).Run()
	if err != nil {
		log.Println(err)
	}
}

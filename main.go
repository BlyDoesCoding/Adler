package main

import (
	"fmt"
	"image/color"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	streamdeck "github.com/magicmonkey/go-streamdeck"
	_ "github.com/magicmonkey/go-streamdeck/devices"
)

var currentpage = []string{"default"}

var globalString string

func main() {

	globalString = "default"

	sd, err := streamdeck.Open()
	if err != nil {
		panic(err)
	}

	currentpage = append(currentpage, getJSONFileNames()...)

	a := app.New()
	w := a.NewWindow("Adler")

	w.Resize(fyne.NewSize(1920/2, 1080/2))
	w.CenterOnScreen()
	x := 5
	y := 3 // Number of columns

	page_2 := container.New(
		layout.NewVBoxLayout(),
	)

	// Create a grid of buttons
	grid := container.NewGridWithColumns(x) // Use x rows
	id := 0
	for j := 0; j < x; j++ {
		for i := 0; i < y; i++ {
			button := widget.NewButton("Button", func(id int) func() {
				return func() {
					pressed(id, a, page_2, sd)
				}
			}(id))
			grid.Add(button)

			id++
		}
	}

	selector := widget.NewSelect(currentpage, func(value string) {

		globalString = value
		updateStreamdeck(sd)

	})

	sd.ButtonPress(func(i int, d *streamdeck.Device, err error) {
		doEvents(i, appendJSONExtension(globalString))

	})

	page_1 := container.New(
		layout.NewVBoxLayout(),
		selector,
		widget.NewButton("New Page", func() {
			dialog.NewEntryDialog("Name of new Page", "Name", func(s string) {

				currentpage = append(currentpage, s)

				selector.SetOptions(currentpage)

				selector.Refresh()
				selector.SetSelected(s)
				globalString = s

				updateStreamdeck(sd)
			}, w).Show()
		}),
		widget.NewLabel(sd.GetName()),

		widget.NewSeparator(),
		grid,
		widget.NewSeparator(),
	)

	page_1.Add(page_2)

	updateStreamdeck(sd)

	selector.SetSelected("default")

	w.SetContent(page_1)

	w.ShowAndRun()

}

func pressed(id int, a fyne.App, page *fyne.Container, sd *streamdeck.Device) {

	addLabelsToPage(a, page, appendJSONExtension(globalString), id, sd)

}

func updateStreamdeck(sd *streamdeck.Device) {

	sd.ClearButtons()

	// Call the function to read all keys and their associated data
	keysAndData, err := readAllKeysAndData(appendJSONExtension(globalString))
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Iterate through the keys and associated data
	for key, _ := range keysAndData {

		sd.WriteColorToButton(key, color.RGBA{255, 255, 255, 255})
	}

}

func addLabelsToPage(a fyne.App, page *fyne.Container, filename string, id int, sd *streamdeck.Device) {
	// Read the data from the JSON file
	data, err := readJSONFromFile(filename, id)
	if err != nil {
		// Handle the error

	}
	var labels []fyne.CanvasObject

	var isInsideAction bool
	var currentAction []string

	for _, item := range data {
		if isInsideAction {
			// Continue adding labels until "STOPPACTION" is encountered
			if item == "STOPPACTION" {
				isInsideAction = false
				// Create a new label with the content of the current action
				label := widget.NewLabel(strings.Join(currentAction, "\n"))
				labels = append(labels, label)
				currentAction = nil
			} else {
				currentAction = append(currentAction, item)
			}
		} else if item == "OPENAPP" {
			isInsideAction = true
			currentAction = append(currentAction, "Open App:")
		} else if item == "OPENSITE" {
			isInsideAction = true
			currentAction = append(currentAction, "Open Site:")
		} else {
			// Create a label with the current string
			label := widget.NewLabel(item)
			labels = append(labels, label)
		}

	}

	// If an action is still open when the loop ends, create a label for it
	if isInsideAction {
		label := widget.NewLabel(strings.Join(currentAction, "\n"))
		labels = append(labels, label)
	}
	page.RemoveAll()

	for _, label := range labels {
		page.Add(label)
	}

	page.Add(widget.NewButton("Add Action", func() {
		addAction(a, id, globalString)
		page.RemoveAll()
		addLabelsToPage(a, page, filename, id, sd)

	}))

	updateStreamdeck(sd)
}

func doEvents(id int, filename string) {
	data, err := readJSONFromFile(filename, id)
	if err != nil {
		fmt.Printf("Error reading JSON from %s: %v\n", filename, err)
		return
	}

	// Join the slice of strings into a single string
	jsonData := strings.Join(data, "")

	// Split the data into chunks whenever "STOPPACTION" is encountered
	chunks := strings.Split(jsonData, "STOPPACTION")

	// Iterate over the chunks and print them
	for _, chunk := range chunks {

		if strings.Contains(chunk, "OPENSITE") {
			site := removeSubstring(chunk, "OPENSITE")
			OpenWebsite(site)

		}
		if strings.Contains(chunk, "OPENAPP") {
			app := removeSubstring(chunk, "OPENAPP")
			startBinary(app)

		}

	}

}

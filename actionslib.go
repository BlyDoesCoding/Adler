package main

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func getNewAction(a fyne.App) []string {
	BinaryPathToApp := widget.NewEntry()
	BinaryPathToApp.SetPlaceHolder("Path To app Binary")

	URL := widget.NewEntry()
	URL.SetPlaceHolder("URL of Website")
	// Create a channel for receiving the result
	resultCh := make(chan []string)

	action := a.NewWindow("Action")
	action.Resize(fyne.NewSize(1920/2, 1080/2))
	action.CenterOnScreen()

	startAppPage := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		BinaryPathToApp,
		layout.NewSpacer(),
		widget.NewButton("Use this action", func() {
			action.Close()
			strList := []string{"OPENAPP", BinaryPathToApp.Text, "STOPPACTION"}

			// Send the result through the channel
			resultCh <- strList
		}),
	)

	OpenWebsite := container.New(
		layout.NewVBoxLayout(),
		layout.NewSpacer(),
		URL,
		layout.NewSpacer(),
		widget.NewButton("Use this action", func() {
			action.Close()
			strList := []string{"OPENSITE", URL.Text, "STOPPACTION"}

			// Send the result through the channel
			resultCh <- strList
		}),
	)

	tabs := container.NewAppTabs(
		container.NewTabItem("Start app", startAppPage),
		container.NewTabItem("Open a Website", OpenWebsite),
	)

	action.SetContent(tabs)
	action.Show()

	// Wait for the result to be sent through the channel
	result := <-resultCh

	return result
}

func addAction(a fyne.App, id int, filename string) {

	err := writeOrUpdateJSONToFile(appendJSONExtension(filename), id, getNewAction(a))
	if err != nil {
		fmt.Println("Error:", err)
	} else {

	}
}

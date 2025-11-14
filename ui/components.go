package ui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
)

func ErrorWindow(app fyne.App, msg string) {
	errorWindow := app.NewWindow("EzPlay")
	errorLabel := widget.NewLabel(msg)
	errorButton := widget.NewButton("Close", func() {
		errorWindow.Close()
	})
	errorWindow.SetContent(container.NewVBox(errorLabel, errorButton))
	errorWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
			errorWindow.Close()
		}
	})
	errorWindow.Show()
}

func NewFileList() *widget.List {
	return widget.NewList(
		func() int { return 0 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {},
	)
}

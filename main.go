package main

import (
	"github.com/kshrs/ezplay/logic"
	"github.com/kshrs/ezplay/ui"
	"fyne.io/fyne/v2/app"
	// "log"
	// "errors"
	// "fmt"
)

func main() {
	app := app.NewWithID("github.com/kshrs/ezplay")
	state := logic.AppState{}

	mainWindow := ui.BuildMainWindow(app, &state)
	mainWindow.Show()
	app.Run()
}

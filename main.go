package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"os"
	// "fmt"
	"log"
	"strings"
)

func returnFileNames(searchString string) ([]string, error) {
	entries, err := os.ReadDir(".")
	if err != nil {
		return nil, err
	}
	var files []string
	for _, entry := range entries {
		if (!entry.IsDir()) {
			if (strings.Contains(entry.Name(), searchString)) {
				files = append(files, entry.Name())	
				
			}
		}
	}
	return files, nil
}

func createFileView() fyne.CanvasObject {
	label := widget.NewLabelWithStyle(
		"Files",
		fyne.TextAlignLeading,
		fyne.TextStyle{Bold: true, Italic: true, Underline: true},
	)
	fileLabel := widget.NewLabel("")

	searchText := widget.NewMultiLineEntry()
	searchText.PlaceHolder = "Enter a text to search"
	submit := widget.NewButton("Search", func() {
		files, err := returnFileNames(searchText.Text)
		if err != nil {
			log.Panic("Error: ", err)
		}
		fileLabel.SetText(strings.Join(files, "\n"))
	})

	return container.NewVBox(label, searchText, submit, fileLabel)

}


func main()  {
	a := app.New()
	w := a.NewWindow("Main Window")

	
	title := widget.NewLabelWithStyle(
		"Welcome To EzPlay!",
		fyne.TextAlignCenter,
		fyne.TextStyle{Bold: true},
	)

	w.SetContent(container.NewVBox(
		title,
		createFileView(),
	))
	w.ShowAndRun()

}

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"log"
	"os"
	"os/exec"
	"strings"
	"errors"
)

var dir string = "/home/kishor/Videos/Data Structures and Algorithms Design/"
func returnFileList(searchString string) ([]string, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var files []string
	for _, entry := range entries {
		if (!entry.IsDir() && strings.Contains(entry.Name(), searchString)) {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func errorWindow(app fyne.App, msg string) {
	w := app.NewWindow("Error Window")
	errorLabel := widget.NewLabel("Error: " + msg)
	closeButton := widget.NewButton("Close", func() {
		w.Close()
	})
	w.SetContent(container.NewVBox(errorLabel,closeButton))
	w.Show()
}

func returnPlayCombo(fileList []string, searchString string, videoFileType string, subFileType string) (string, string, error) {
	var videoFileName string = ""
	var subFileName string = ""

	for _, file := range fileList {
		if (strings.Contains(file, searchString) && strings.Contains(file, videoFileType)) {
			videoFileName = file
		}
		if (strings.Contains(file, searchString) && strings.Contains(file, subFileType)) {
			subFileName = file
		}
	}
	if (videoFileName == "" && subFileName == "") {
		return "", "", errors.New("No Video and Subtrack Found")
	}
	if (subFileName == "") {
		return videoFileName, "", nil
	}

	return strings.Join([]string{dir, videoFileName}, ""),
			strings.Join([]string{dir, subFileName}, ""),
			nil
}

func main() {
	a := app.New()
	w := a.NewWindow("Main Window")


	fileList, err := returnFileList("")
	if err != nil {
		errorWindow(a, err.Error())
		return
	}
	fileListLabel := widget.NewLabel(strings.Join(fileList, "\n"))

	searchEntry := widget.NewMultiLineEntry()
	searchEntry.PlaceHolder = "Type to search files"
	searchButton := widget.NewButton("Search", func() {
		fileList, err := returnFileList(searchEntry.Text)
		if err != nil {
			errorWindow(a, err.Error())
			return
		}
		fileListLabel.SetText(strings.Join(fileList, "\n"))
	})

	playButton := widget.NewButton("Play", func() {
		videoFileName, subFileName, err := returnPlayCombo(fileList, searchEntry.Text, ".mp4", ".srt")
		if err != nil {
			errorWindow(a, err.Error())
			return
		}

		playCmd := exec.Command("mpv", videoFileName, strings.Join([]string{"--sub-file", subFileName}, "="))
		log.Println(playCmd)
		err = playCmd.Run()
		if err != nil {
			errorWindow(a, "MPV Error")
			return
		}


	})
	

	w.SetContent(container.NewVBox(searchEntry, searchButton, fileListLabel, playButton))
	w.ShowAndRun()
}

package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"

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
	w := app.NewWindow("EzPlay")
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
		return "", "", errors.New("No Video and Subtrack Found or Selected")
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
	w := a.NewWindow("EzPlay")


	fileList, err := returnFileList("")
	if err != nil {
		errorWindow(a, err.Error())
		return
	}
	fileListLabel := widget.NewLabel(strings.Join(fileList, "\n"))
	scrollFileList := container.NewVScroll(fileListLabel)
	// scrollFileList.Resize(fyne.NewSize(500, 150))
	fixedScrollFileList := container.New(
		layout.NewGridWrapLayout(fyne.NewSize(500, 150)),
		scrollFileList,
	)


	// searchEntry := widget.NewMultiLineEntry()
	searchEntry := widget.NewEntry()
	searchEntry.PlaceHolder = "Type to search files"
	searchButton := widget.NewButton("Search", func() {
		fileList, err := returnFileList(searchEntry.Text)
		if err != nil {
			errorWindow(a, err.Error())
			return
		}
		fileListLabel.SetText(strings.Join(fileList, "\n"))
	})

	searchEntry.OnSubmitted = func(_ string) {
		searchButton.Tapped(nil)
	}


	
	videoType := ".mp4"
	videoTypes := []string{".mp4", ".mp3", ".mkv"}
	comboVideoType := widget.NewSelect(videoTypes, func(value string) {videoType = value})
	comboVideoType.PlaceHolder = ".mp4"

	subType := ".srt"
	subTypes := []string{"None", ".srt"}
	comboSubType := widget.NewSelect(subTypes, func(value string) {
		if (value == "None") {
			subType = ""
		} else {
			subType = value
		}
	})
	comboSubType.PlaceHolder = ".srt"

	playButton := widget.NewButton("Play", func() {
		videoFileName, subFileName, err := returnPlayCombo(fileList, searchEntry.Text, videoType, subType)
		if err != nil {
			errorWindow(a, err.Error())
			return
		}

		var playCmd *exec.Cmd
		if (subType == "None") {
			playCmd = exec.Command("mpv", videoFileName)
		} else {
			playCmd = exec.Command("mpv", videoFileName, strings.Join([]string{"--sub-file", subFileName}, "="))
		}
		log.Println(playCmd)
		err = playCmd.Run()
		if err != nil {
			errorWindow(a, "MPV Error")
			return
		}
	})

	

	w.SetContent(
		container.NewVBox(
			searchEntry, 
			searchButton,
			// scrollFileList,
			fixedScrollFileList,
			layout.NewSpacer(),
			container.NewHBox(comboVideoType, comboSubType),
			playButton,
	))
	w.Resize(fyne.NewSize(500, 300))
	w.Show()
	
	fyne.Do(func() {
		w.Canvas().Focus(searchEntry)
	})

	a.Run()
}

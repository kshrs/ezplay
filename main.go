package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/dialog"

	// "log"
	"os"
	"os/exec"
	"strings"
	"errors"
	// "fmt"
	"path/filepath"
)

func returnFileList(searchString string, folderPath string) ([]string, error) {
	if folderPath == "" {
		fileList := []string{"Select a Folder to Search"}
		return fileList, nil
	}
	entries, err := os.ReadDir(folderPath)
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

	w.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if (ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter) {
			closeButton.Tapped(nil)
		}

	})
	w.Show()
}

func returnPlayCombo(fileList []string, searchString string, videoFileType string, subFileType string, folderPath string) (string, string, error) {
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

	return filepath.Join(folderPath, videoFileName),
			filepath.Join(folderPath, subFileName),
			nil
}

func main() {
	a := app.NewWithID("github.com/kshrs/ezplay")
	w := a.NewWindow("EzPlay")

	folderPath := ""


	fileList, err := returnFileList("", folderPath)
	if err != nil {
		errorWindow(a, err.Error())
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
	// searchButtonFunction := func() {
	// 	fileList, err := returnFileList(searchEntry.Text, folderPath)
	// 	if err != nil {
	// 		errorWindow(a, err.Error())
	// 		return
	// 	}
	// 	fileListLabel.SetText(strings.Join(fileList, "\n"))
	// }
	// searchButton := widget.NewButton("Search", searchButtonFunction)
	//
	// searchEntry.OnSubmitted = func(_ string) {
	// 	searchButton.Tapped(nil)
	// }
	searchEntry.OnChanged = func(_ string) {
		fileList, err := returnFileList(searchEntry.Text, folderPath)
		if err != nil {
			errorWindow(a, err.Error())
			return
		}
		fileListLabel.SetText(strings.Join(fileList, "\n"))

	}




	folderLabel := widget.NewLabel("No Directory Selected")
	folderLabel.Wrapping = fyne.TextWrapWord
	chooseFolderButton := widget.NewButton("Choose Directory", func() {
		fd := dialog.NewFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				errorWindow(a, err.Error())
				return
			}
			if folder == nil {
				return
			}
			folderPath = folder.Path()
			if (folderPath != "") {
				folderLabel.SetText("Path: " + folderPath)
			}
			fileList, err = returnFileList(searchEntry.Text, folderPath)
			if err != nil {
				errorWindow(a, err.Error())
				return
			}
			fileListLabel.SetText(strings.Join(fileList, "\n"))
		}, w)

		fd.Show()
	})



	
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

	playerType := "mpv"
	playerTypes := []string{"mpv", "vlc"}
	comboPlayerType := widget.NewSelect(playerTypes, func(value string) {
		playerType = value
	})
	comboPlayerType.PlaceHolder = "mpv"

	playButton := widget.NewButton("Play", func() {
		videoFileName, subFileName, err := returnPlayCombo(fileList, searchEntry.Text, videoType, subType, folderPath)
		if err != nil {
			errorWindow(a, err.Error())
			return
		}

		var playCmd *exec.Cmd
		if (subType == "None") {
			playCmd = exec.Command(playerType, videoFileName)
		} else {
			playCmd = exec.Command(playerType, videoFileName, strings.Join([]string{"--sub-file", subFileName}, "="))
		}
		err = playCmd.Run()
		if err != nil {
			errorWindow(a, strings.Join([]string{playerType, err.Error()}, " "))
			return
		}
	})
	searchEntry.OnSubmitted = func(_ string) {
		if folderPath == "" {
			chooseFolderButton.Tapped(nil)
		} else {
			playButton.Tapped(nil)
		}
	}

	

	w.SetContent(
		container.NewVBox(
			// container.NewHBox(folderLabel, layout.NewSpacer(), chooseFolderButton),
			// container.NewBorder Doesn't require layout.NewSpacer() and works well with fyne.TextWrapWord -> label.Wrapping
			container.NewBorder(nil, nil, nil, chooseFolderButton, folderLabel),
			searchEntry, 
			// searchButton,
			// scrollFileList,
			fixedScrollFileList,
			layout.NewSpacer(),
			container.NewHBox(comboVideoType, comboSubType, comboPlayerType),
			playButton,
	))
	w.Resize(fyne.NewSize(500, 300))
	w.Show()
	
	fyne.Do(func() {
		w.Canvas().Focus(searchEntry)
	})

	a.Run()
}

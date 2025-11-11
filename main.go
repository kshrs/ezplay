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
	// "errors"
	// "fmt"
	"path/filepath"
)

// Function: searchFiles returns a slice of string which containes the search results 
// for a specific search string and an error message.
func searchFiles(searchString string, dirToSearch string) ([]string, error) {
	if dirToSearch == "" {
		return []string{"Select a Folder to Search"}, nil
	}
	var files []string

	entries, err := os.ReadDir(dirToSearch)
	if err != nil {
		return nil, err
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			if strings.Contains(entry.Name(), searchString) {
				files = append(files, entry.Name())
			}
		}
	}
	return files, nil
}

// Function: errorWindow creates a new popup window to display the error message.
func errorWindow(app fyne.App, msg string) {
	errWindow := app.NewWindow("EzPlay")
	errLabel := widget.NewLabel("Error: " + msg)
	closeButton := widget.NewButton("Close", func() {
		errWindow.Close()
	})

	errWindow.SetContent(container.NewVBox(errLabel, closeButton))
	errWindow.Show()

	errWindow.Canvas().SetOnTypedKey(func(ev *fyne.KeyEvent) {
		if(ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter) {
			closeButton.Tapped(nil)
		}
	})
}

// PlayFile struct is a Result struct which holds the absolute location of the video and sub file as string
type PlayFile struct {
	VideoFileLocation string
	SubFileLocation string
}

// Function: locatePlayFiles return a *PlayFile struct which contains the locations of the video,sub file and an error message
func locatePlayFiles(files []string, searchString string, videoFileType string, subFileType string, fileDir string) (*PlayFile, error) {
	var (
		videoFile string
		subFile string
	)

	for _, file := range files {
		if !strings.Contains(file, searchString) {
			continue
		}

		if videoFile == "" && strings.Contains(file, videoFileType) {
			videoFile = file
		} else if subFile == "" && strings.Contains(file, subFileType) {
			subFile = file
		}
		
		if videoFile != "" && subFile != "" {
			break
		}
	}

	return &PlayFile{
		VideoFileLocation: filepath.Join(fileDir, videoFile),
		SubFileLocation: filepath.Join(fileDir, subFile),
	}, nil
}

// func searchFileWidget(app fyne.App, dirOfFiles string) *fyne.Container {
//
// 	searchEntry := widget.NewEntry()
// 	searchEntry.PlaceHolder = "Type to search files"
//
// 	searchEntry.OnChanged = func(_ string) {
// 		fileList, err := searchFiles(searchEntry.Text, dirOfFiles)
// 		if err != nil {
// 			errorWindow(app, err.Error())
// 			return
// 		}
// 		fileListLabel.SetText(strings.Join(fileList, "\n"))
//
// 	}
//
// 	files, err := searchFiles("", dirOfFiles)
// 	if err != nil {
// 		errorWindow(app, err.Error())
// 	}
// 	filesDisplayLabel := widget.NewLabel(strings.Join(files, "\n"))
// 	fileDisplayScrollContainer := container.NewVScroll(filesDisplayLabel)
//
// 	return container.NewVBox(
// 		container.New(
// 			layout.NewGridWrapLayout(fyne.NewSize(500, 150)),
// 			fileDisplayScrollContainer,
// 		),
// 	)
// }
//


func main() {
	app := app.NewWithID("github.com/kshrs/ezplay")
	mainWindow := app.NewWindow("EzPlay")

	var dirOfFiles string

	files, err := searchFiles("", dirOfFiles)
	if err != nil {
		errorWindow(app, err.Error())
	}
	fileListLabel := widget.NewLabel(strings.Join(files, "\n"))
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
		fileList, err := searchFiles(searchEntry.Text, dirOfFiles)
		if err != nil {
			errorWindow(app, err.Error())
			return
		}
		fileListLabel.SetText(strings.Join(fileList, "\n"))

	}




	folderLabel := widget.NewLabel("No Directory Selected")
	folderLabel.Wrapping = fyne.TextWrapWord
	chooseFolderButton := widget.NewButton("Choose Directory", func() {
		fd := dialog.NewFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				errorWindow(app, err.Error())
				return
			}
			if folder == nil {
				return
			}
			dirOfFiles = folder.Path()
			if (dirOfFiles != "") {
				folderLabel.SetText("Path: " + dirOfFiles)
			}
			files, err = searchFiles(searchEntry.Text, dirOfFiles)
			if err != nil {
				errorWindow(app, err.Error())
				return
			}
			fileListLabel.SetText(strings.Join(files, "\n"))
		}, mainWindow)

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
		playFile, err := locatePlayFiles(files, searchEntry.Text, videoType, subType, dirOfFiles)
		if err != nil {
			errorWindow(app, err.Error())
			return
		}

		var playCmd *exec.Cmd
		if (subType == "None") {
			playCmd = exec.Command(playerType, playFile.VideoFileLocation)
		} else {
			playCmd = exec.Command(playerType, playFile.VideoFileLocation, strings.Join([]string{"--sub-file", playFile.SubFileLocation}, "="))
		}
		err = playCmd.Run()
		if err != nil {
			errorWindow(app, strings.Join([]string{playerType, err.Error()}, " "))
			return
		}
	})
	searchEntry.OnSubmitted = func(_ string) {
		if dirOfFiles == "" {
			chooseFolderButton.Tapped(nil)
		} else {
			playButton.Tapped(nil)
		}
	}

	

	mainWindow.SetContent(
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
	mainWindow.Resize(fyne.NewSize(500, 300))
	mainWindow.Show()
	
	fyne.Do(func() {
		mainWindow.Canvas().Focus(searchEntry)
	})

	app.Run()
}

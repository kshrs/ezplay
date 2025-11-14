package ui

import (
	"github.com/kshrs/ezplay/logic"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"strings"

	"os/exec"
)

func BuildMainWindow(app fyne.App, state *logic.AppState) fyne.Window {

	var playOptions logic.PlayOptions = logic.PlayOptions{
		PlayerType: "mpv",
		VideoType: ".mp4",
		SubType: ".srt",
	}

	mainWindow := app.NewWindow("EzPlay")

	fileList := widget.NewList(
		func() int { return len(state.Filtered) },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(i widget.ListItemID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(state.Filtered[i])
		},
	)

	folderLabel := widget.NewLabel("No Directory Selected")
	folderLabel.Wrapping = fyne.TextWrapWord
	chooseFolderButton := widget.NewButton("Choose Directory", func() {
		fd := dialog.NewFolderOpen(func(folder fyne.ListableURI, err error) {
			if err != nil {
				ErrorWindow(app, err.Error())
				return
			}
			if folder == nil {
				return
			}
			state.Dir = folder.Path()
			if (state.Dir != "") {
				folderLabel.SetText("Path: " + state.Dir)
			}
			state.Files, err = logic.ListFiles(state.Dir)
			state.Filtered = state.Files
			fileList.Refresh()
			if err != nil {
				ErrorWindow(app, err.Error())
				return
			}
		}, mainWindow)

		fd.Show()
	})

	entry := widget.NewEntry()
	entry.SetPlaceHolder("Type to search files")
	

	entry.OnChanged = func(text string) {
		state.Filtered = logic.FilterFiles(state.Files, text)
		fileList.Refresh()
	}

	videoTypes := []string{".mp4", ".mp3", ".mkv"}
	comboVideoType := widget.NewSelect(videoTypes, func(value string) {playOptions.VideoType = value})
	comboVideoType.PlaceHolder = ".mp4"

	subTypes := []string{"None", ".srt"}
	comboSubType := widget.NewSelect(subTypes, func(value string) {
		playOptions.SubType = value
	})
	comboSubType.PlaceHolder = ".srt"

	playerType := "mpv"
	playerTypes := []string{"mpv", "vlc"}
	comboPlayerType := widget.NewSelect(playerTypes, func(value string) {
		playerType = value
	})
	comboPlayerType.PlaceHolder = "mpv"

	playButton := widget.NewButton("Play", func() {
		playFile, err := logic.LocatePlayFiles(state.Files, entry.Text, playOptions.VideoType, playOptions.SubType, state.Dir)
		if err != nil {
			ErrorWindow(app, err.Error())
			return
		}

		var playCmd *exec.Cmd
		if (playOptions.SubType == "None") {
			playCmd = exec.Command(playerType, playFile.VideoFileLocation)
		} else {
			playCmd = exec.Command(playerType, playFile.VideoFileLocation, strings.Join([]string{"--sub-file", playFile.SubFileLocation}, "="))
		}
		err = playCmd.Run()
		if err != nil {
			ErrorWindow(app, strings.Join([]string{playerType, err.Error()}, " "))
			return
		}
	})

	entry.OnSubmitted = func(_ string) {
		if state.Dir == "" {
			chooseFolderButton.Tapped(nil)
		} else {
			playButton.Tapped(nil)
		}
	}


	mainWindow.SetContent(
		container.NewVBox(
			container.NewBorder(nil, nil, nil, chooseFolderButton, folderLabel),
			entry,
			layout.NewSpacer(),
			container.New(layout.NewGridWrapLayout(fyne.NewSize(500, 150)), fileList),
			container.NewHBox(comboVideoType, comboSubType, comboPlayerType),
			playButton,
		),	
	)

	mainWindow.Resize(fyne.NewSize(500, 300))
	fyne.Do(func() {
		mainWindow.Canvas().Focus(entry)
	})
	return mainWindow
}



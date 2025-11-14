package logic

import (
	"os"
	"strings"
	"path/filepath"
)

// PlayFile struct is a Result struct which holds the absolute location of the video and sub file as string
type PlayFile struct {
	VideoFileLocation string
	SubFileLocation string
}
type PlayOptions struct {
	PlayerType string
	VideoType string
	SubType string
}


func ListFiles(dir string) ([]string, error) {
	if dir == "" {
		return []string{"Select a Directory to view them"}, nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	var files []string

	for _, entry := range entries {
		if !entry.IsDir() {
			files = append(files, entry.Name())
		}
	}
	return files, nil
}

func FilterFiles(files []string, query string) []string {
	if query == "" {
		return files
	}
	var result []string
	for _, file := range files {
		if strings.Contains(strings.ToLower(file), strings.ToLower(query)) {
			result = append(result, file)
		}
	}
	return result
}

// Function: LocatePlayFiles return a *PlayFile struct which contains the locations of the video,sub file and an error message
func LocatePlayFiles(files []string, searchString string, videoFileType string, subFileType string, fileDir string) (*PlayFile, error) {
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

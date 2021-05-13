package main

import (
	"fmt"

	"github.com/grumpypixel/go-webget"
)

func main() {
	url := "https://upload.wikimedia.org/wikipedia/commons/d/d6/Wp-w4-big.jpg"
	targetDir := "."
	targetFilename := "" // use original filename
	options := webget.Options{
		ProgressHandler: MyProgress{},
	}
	webget.Download(url, targetDir, targetFilename, &options)
}

type MyProgress struct{}

func (p MyProgress) Start(sourceURL, filename string) {
	fmt.Printf("Starting download %s\n", filename)
}

func (p MyProgress) Update(sourceURL, filename string, percentage float64, bytes, size uint64) {
	fmt.Printf("\rDownloading %s: %v bytes [%.2f%%]", filename, bytes, percentage)
}

func (p MyProgress) Done(sourceURL, filename, targetFilePath string) {
	fmt.Printf("\nFinished download %s\n", filename)
}

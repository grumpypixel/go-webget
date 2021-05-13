package main

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/grumpypixel/go-webget"
)

func main() {
	if err := downloadToFile(); err != nil {
		fmt.Println(err)
	}
	if err := downloadToBuffer(); err != nil {
		fmt.Println(err)
	}
}

func downloadToFile() error {
	fmt.Println("Downloading to file...")
	url := "https://upload.wikimedia.org/wikipedia/commons/d/d6/Wp-w4-big.jpg"
	targetDir := "./"
	targetFilename := "" // use original filename
	options := webget.Options{
		ProgressHandler: MyProgress{},
		Timeout:         time.Second * 60,
	}
	return webget.DownloadToFile(url, targetDir, targetFilename, &options)
}

func downloadToBuffer() error {
	fmt.Println("\nDownloading to buffer and save to file...")
	url := "https://golang.org/doc/gopher/appenginegophercolor.jpg"
	options := webget.Options{
		ProgressHandler: MyProgress{},
		Timeout:         time.Second * 60,
	}
	bytes, err := webget.DownloadToBuffer(url, &options)
	if err != nil {
		return err
	}
	filename := path.Base(url)
	if _, err := writeBufferToFile(filename, bytes); err != nil {
		return err
	}
	return nil
}

func writeBufferToFile(filepath string, bytes []byte) (int, error) {
	file, err := os.Create(filepath)
	if err != nil {
		return 0, err
	}
	writer := bufio.NewWriter(file)
	nn, err := writer.Write(bytes)
	if err != nil {
		return nn, nil
	}
	writer.Flush()
	return nn, nil
}

type MyProgress struct{}

func (p MyProgress) Start(sourceURL string) {
	fmt.Printf("Starting %s\n", sourceURL)
}

func (p MyProgress) Update(sourceURL string, percentage float64, bytesRead, contentLength int64) {
	fmt.Printf("\rDownloading %s: %v bytes [%.2f%%]", path.Base(sourceURL), bytesRead, percentage)
}

func (p MyProgress) Done(sourceURL string) {
	fmt.Printf("\nDone %s\n", path.Base(sourceURL))
}

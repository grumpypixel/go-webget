package main

import (
	"flag"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/grumpypixel/go-webget"
)

func main() {
	url, targetDir, timeout := parseParams()
	if url == "" {
		fmt.Println("Given URL is empty. Bye.")
		return
	}

	targetFilename := "" // use original filename
	options := webget.Options{
		ProgressHandler: MyProgress{},
		Timeout:         timeout,
		CreateTargetDir: true,
	}
	if err := webget.DownloadToFile(url, targetDir, targetFilename, &options); err != nil {
		fmt.Println(err)
	}
}

func parseParams() (string, string, time.Duration) {
	url := flag.String("url", "", "")
	targetDir := flag.String("target", "./", "")
	timeout := flag.Duration("timeout", time.Second*60, "")
	flag.Parse()

	if *url == "" {
		args := os.Args[1:]
		if len(args) > 0 {
			*url = args[0]
		}
	}
	return *url, *targetDir, *timeout
}

type MyProgress struct{}

func (p MyProgress) Start(sourceURL string) {
}

func (p MyProgress) Update(sourceURL string, percentage float64, bytesRead, contentLength int64) {
	fmt.Printf("\rDownloading %s: %v bytes [%.2f%%]", path.Base(sourceURL), bytesRead, percentage)
}

func (p MyProgress) Done(sourceURL string) {
	fmt.Printf("\n")
}

package main

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/grumpypixel/go-webget"
)

func main() {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)

	urls := []string{
		"https://photojournal.jpl.nasa.gov/jpeg/PIA22358.jpg",
		"https://photojournal.jpl.nasa.gov/jpeg/PIA23126.jpg",
		"https://photojournal.jpl.nasa.gov/jpeg/PIA23127.jpg",
		"https://photojournal.jpl.nasa.gov/jpeg/PIA23405.jpg",
		"https://photojournal.jpl.nasa.gov/jpeg/PIA23646.jpg",
		"https://photojournal.jpl.nasa.gov/jpeg/PIA23647.jpg",
		"https://photojournal.jpl.nasa.gov/jpeg/PIA24472.jpg",
	}

	targetDir := "./downloads"
	targetFilename := "" // use original filename

	waitGroup := sync.WaitGroup{}
	waitGroup.Add(len(urls))
	for i, url := range urls {
		options := webget.Options{
			ProgressHandler: MyProgress{index: i, waitGroup: &waitGroup},
			CreateTargetDir: true,
		}
		go webget.Download(url, targetDir, targetFilename, &options)
	}
	waitGroup.Wait()
}

type MyProgress struct {
	index     int
	waitGroup *sync.WaitGroup
}

func (p MyProgress) Start(sourceURL, filename string) {
	fmt.Printf("Starting download %s (#%d)\n", filename, p.index)
}

func (p MyProgress) Update(sourceURL, filename string, percentage float64, bytes, size uint64) {
	fmt.Printf(".")
}

func (p MyProgress) Done(sourceURL, filename, targetFilePath string) {
	fmt.Printf("\nFinished download %s (#%d)\n", filename, p.index)
	p.waitGroup.Done()
}

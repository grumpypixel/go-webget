package main

import (
	"fmt"
	"path"
	"runtime"
	"sync"
	"time"

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
	for i, url := range urls {
		options := webget.Options{
			ProgressHandler: MyProgress{index: i, waitGroup: &waitGroup},
			Timeout:         time.Second * 60,
			CreateTargetDir: true,
		}
		waitGroup.Add(1)
		go webget.DownloadToFile(url, targetDir, targetFilename, &options)
	}
	waitGroup.Wait()
}

type MyProgress struct {
	index     int
	waitGroup *sync.WaitGroup
}

func (p MyProgress) Start(sourceURL string) {
	fmt.Printf("Starting %s [#%d]\n", path.Base(sourceURL), p.index)
}

func (p MyProgress) Update(sourceURL string, percentage float64, bytesRead, contentLength int64) {
	fmt.Printf(".")
}

func (p MyProgress) Done(sourceURL string) {
	fmt.Printf("\nDone %s [#%d]\n", path.Base(sourceURL), p.index)
	p.waitGroup.Done()
}

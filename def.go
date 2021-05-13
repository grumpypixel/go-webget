package webget

import "time"

const (
	DefaultTimeout = time.Second * 60
	TempExtension  = ".godownload"
)

type ProgressHandler interface {
	Start(soureURL string)
	Update(soureURL string, percentage float64, bytes, contentLength int64)
	Done(soureURL string)
}

type Options struct {
	ProgressHandler ProgressHandler
	Timeout         time.Duration
	CreateTargetDir bool
}

func DefaultOptions() *Options {
	return &Options{
		ProgressHandler: nil,
		Timeout:         DefaultTimeout,
		CreateTargetDir: false,
	}
}

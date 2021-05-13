package webget

const (
	tempExtension = ".godownload"
)

type ProgressHandler interface {
	Start(soureURL, filename string)
	Update(soureURL, filename string, percentage float64, bytes, size uint64)
	Done(soureURL, filename, targetFilePath string)
}

type Options struct {
	ProgressHandler ProgressHandler
	CreateTargetDir bool
}

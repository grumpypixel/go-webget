package webget

const (
	tempExtension = ".godownload"
)

type ProgressHandler interface {
	Start(soureURL string)
	Update(soureURL string, percentage float64, bytes, contentLength int64)
	Done(soureURL string)
}

type Options struct {
	ProgressHandler ProgressHandler
	CreateTargetDir bool
}

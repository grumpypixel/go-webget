package webget

type ProgressWriter struct {
	URL             string
	Filename        string
	BytesRead       int64
	ContentLength   int64
	ProgressHandler *ProgressHandler
	CustomPayload   interface{}
}

func (fp *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	fp.BytesRead += int64(n)
	fp.Update()
	return n, nil
}

func (fp *ProgressWriter) Update() {
	if fp.ProgressHandler != nil {
		var percentage float64
		if fp.ContentLength > 0 {
			percentage = float64(fp.BytesRead) / float64(fp.ContentLength) * 100.0
		} else {
			percentage = -1.0
		}
		(*fp.ProgressHandler).Update(fp.URL, percentage, fp.BytesRead, fp.ContentLength)
	}
}

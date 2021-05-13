package webget

type ProgressWriter struct {
	URL             string
	Filename        string
	Bytes           uint64
	Size            uint64
	ProgressHandler *ProgressHandler
	CustomPayload   interface{}
}

func (fp *ProgressWriter) Write(p []byte) (int, error) {
	n := len(p)
	fp.Bytes += uint64(n)
	fp.Update()
	return n, nil
}

func (fp *ProgressWriter) Update() {
	if fp.ProgressHandler != nil {
		var percentage float64
		if fp.Size > 0 {
			percentage = float64(fp.Bytes) / float64(fp.Size) * 100.0
		} else {
			percentage = -1.0
		}
		(*fp.ProgressHandler).Update(fp.URL, fp.Filename, percentage, fp.Bytes, fp.Size)
	}
}

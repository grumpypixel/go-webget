package webget

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
)

func DownloadToFile(sourceURL, targetDir, targetFilename string, options *Options) error {
	if options == nil {
		options = &Options{}
	}

	if options.CreateTargetDir {
		if err := os.MkdirAll(targetDir, os.ModePerm); err != nil {
			return err
		}
	}

	if _, err := os.Stat(targetDir); os.IsNotExist(err) {
		return err
	}

	var filename string
	if targetFilename == "" {
		filename = path.Base(sourceURL)
	} else {
		filename = targetFilename
	}

	targetFilePath := path.Join(targetDir, filename)
	tempFilePath := targetFilePath + tempExtension

	if options.ProgressHandler != nil {
		options.ProgressHandler.Start(sourceURL)
	}

	if err := downloadFile(sourceURL, tempFilePath, filename, options.ProgressHandler); err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, targetFilePath); err != nil {
		return err
	}

	if options.ProgressHandler != nil {
		options.ProgressHandler.Done(sourceURL)
	}
	return nil
}

func DownloadToBuffer(sourceURL string, options *Options) ([]byte, error) {
	if options == nil {
		options = &Options{}
	}

	if options.ProgressHandler != nil {
		options.ProgressHandler.Start(sourceURL)
	}

	filename := path.Base(sourceURL)
	bytes, err := downloadToBuffer(sourceURL, filename, options.ProgressHandler)
	if err != nil {
		return nil, err
	}

	if options.ProgressHandler != nil {
		options.ProgressHandler.Done(sourceURL)
	}
	return bytes, nil
}

func downloadFile(sourceURL, filepath, filename string, progressHandler ProgressHandler) error {
	resp, err := httpGet(sourceURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	var handler *ProgressHandler
	if progressHandler != nil {
		handler = &progressHandler
	}

	fileProgress := &ProgressWriter{URL: sourceURL, Filename: filename, ProgressHandler: handler, ContentLength: resp.ContentLength}
	if _, err = io.Copy(file, io.TeeReader(resp.Body, fileProgress)); err != nil {
		return err
	}
	return nil
}

func downloadToBuffer(sourceURL, filename string, progressHandler ProgressHandler) ([]byte, error) {
	resp, err := httpGet(sourceURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var handler *ProgressHandler
	if progressHandler != nil {
		handler = &progressHandler
	}

	contentLength := resp.ContentLength

	var buf *bytes.Buffer
	if contentLength > 0 {
		buf = bytes.NewBuffer(make([]byte, 0, contentLength))
	} else {
		buf = bytes.NewBuffer(make([]byte, 0))
	}

	fileProgress := &ProgressWriter{URL: sourceURL, Filename: filename, ProgressHandler: handler, ContentLength: contentLength}
	if _, err = io.Copy(buf, io.TeeReader(resp.Body, fileProgress)); err != nil {
		return nil, err
	}
	return buf.Bytes(), err
}

func httpGet(sourceURL string) (*http.Response, error) {
	resp, err := http.Get(sourceURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 && resp.StatusCode != 206 {
		return nil, fmt.Errorf(fmt.Sprintf("http error: status code %d", resp.StatusCode))
	}
	return resp, nil
}

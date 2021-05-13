package webget

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"time"
)

func DownloadToFile(sourceURL, targetDir, targetFilename string, options *Options) error {
	options = validateOptions(options)
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
	tempFilePath := targetFilePath + TempExtension

	if options.ProgressHandler != nil {
		options.ProgressHandler.Start(sourceURL)
	}

	if err := downloadFile(sourceURL, tempFilePath, filename, options.Timeout, options.ProgressHandler); err != nil {
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
	options = validateOptions(options)

	if options.ProgressHandler != nil {
		options.ProgressHandler.Start(sourceURL)
	}

	filename := path.Base(sourceURL)
	bytes, err := downloadToBuffer(sourceURL, filename, options.Timeout, options.ProgressHandler)
	if err != nil {
		return nil, err
	}

	if options.ProgressHandler != nil {
		options.ProgressHandler.Done(sourceURL)
	}
	return bytes, nil
}

func downloadFile(sourceURL, filepath, filename string, timeout time.Duration, progressHandler ProgressHandler) error {
	resp, _, err := httpGet(sourceURL, timeout)
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

func downloadToBuffer(sourceURL, filename string, timeout time.Duration, progressHandler ProgressHandler) ([]byte, error) {
	resp, _, err := httpGet(sourceURL, timeout)
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

func httpGet(sourceURL string, timeout time.Duration) (*http.Response, context.CancelFunc, error) {
	request, err := http.NewRequest("GET", sourceURL, nil)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancelFunc := context.WithTimeout(context.Background(), timeout)
	request = request.WithContext(ctx)

	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		cancelFunc()
		return nil, nil, err
	}

	if resp.StatusCode != 200 {
		cancelFunc()
		defer resp.Body.Close()
		return nil, nil, fmt.Errorf("invalid response; status code: %s", resp.Status)
	}

	return resp, cancelFunc, nil
}

func validateOptions(options *Options) *Options {
	if options != nil {
		if options.Timeout == 0 {
			options.Timeout = DefaultTimeout
		}
		return options
	}
	return DefaultOptions()
}

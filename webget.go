package webget

import (
	"io"
	"net/http"
	"os"
	"path"
)

func Download(sourceURL, targetDir, targetFilename string, options *Options) error {
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
		options.ProgressHandler.Start(sourceURL, filename)
	}

	if err := downloadFile(sourceURL, tempFilePath, filename, options.ProgressHandler); err != nil {
		return err
	}

	if err := os.Rename(tempFilePath, targetFilePath); err != nil {
		return err
	}

	if options.ProgressHandler != nil {
		options.ProgressHandler.Done(sourceURL, filename, targetFilePath)
	}
	return nil
}

func downloadFile(sourceURL, filepath, filename string, progressHandler ProgressHandler) error {
	file, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer file.Close()

	resp, err := http.Get(sourceURL)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var handler *ProgressHandler
	if progressHandler != nil {
		handler = &progressHandler
	}

	fileProgress := &ProgressWriter{URL: sourceURL, Filename: filename, ProgressHandler: handler, Size: uint64(resp.ContentLength)}
	if _, err = io.Copy(file, io.TeeReader(resp.Body, fileProgress)); err != nil {
		return err
	}
	return nil
}

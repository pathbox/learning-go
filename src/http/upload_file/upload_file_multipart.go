package main

import (
	"bytes"
	"io"
	"mime/multipart"
	"net"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"
)

func main() {
	fileDir, _ := os.Getwd()
	fileName := "upload-file.txt"
	filePath := path.Join(fileDir, fileName)

	file, _ := os.Open(filePath)
	file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("file", filepath.Base(file.Name()))
	io.Copy(part, file)
	// If you don't close it, your request will be missing the terminating boundary.
	// can not use `defer writer.Close()`
	// Close() before do request
	writer.Close()

	r, _ := http.NewRequest("POST", "http://example.com", body)
	r.Header.Add("Content-Type", writer.FormDataContentType()) // It must do
	client := NewHTTPClient()
	client.Do(r)
}

func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 30 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   15 * time.Second,
		ResponseHeaderTimeout: 30 * time.Second,
	}

	client := &http.Client{
		Timeout:   30 * time.Second,
		Transport: transport,
	}
	return client
}

package main

import (
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"time"
)

func main() {
	client := NewHTTPClient()
	filepath := "pdf_file.pdf" // 可以换成 html等文件
	url := "http://www.baidu.com"

	outFile, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	fmt.Println("response header: ", resp.Header)
	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	fmt.Println("response body length: ", len(dump))

	// Write the body to file
	_, err = io.Copy(outFile, resp.Body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Download file done")
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

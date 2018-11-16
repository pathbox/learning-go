package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"os"
	"time"
)

func main() {

	filepath := "pdf_file.pdf" // 可以换成 html等文件
	url := "http://www.baidu.com"

	outFile, err := os.Create(filepath)
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	body, err := sendHTTPRequest(url, nil, 20)
	if err != nil {
		panic(err)
	}

	n, err := outFile.Write(body)
	if err != nil {
		panic(err)
	}
	fmt.Println("Write file size: ", n)
	// client := NewHTTPClient()
	// req, err := http.NewRequest("GET", url, nil)
	// if err != nil {
	// 	panic(err)
	// }

	// resp, err := client.Do(req)
	// if err != nil {
	// 	panic(err)
	// }
	// defer resp.Body.Close()
	// fmt.Println("response header: ", resp.Header)
	// dump, err := httputil.DumpResponse(resp, true)
	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println("response body length: ", len(dump))

	// // Write the body to file
	// _, err = io.Copy(outFile, resp.Body)
	// if err != nil {
	// 	panic(err)
	// }
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

func sendHTTPRequest(urlPath string, params map[string]interface{}, timeOut uint32) (res []byte, err error) {
	reqURL, err := url.Parse(urlPath)
	if err != nil {
		return
	}

	reqParams := reqURL.Query()
	for k, v := range params {
		reqParams.Set(k, v.(string))
	}
	reqURL.RawQuery = reqParams.Encode()

	// 设置超时，如果为0,则不超时
	client := newTimeoutHTTPClient(time.Duration(timeOut) * time.Second)
	result, err := client.Get(reqURL.String())
	if err != nil {
		return
	}
	defer result.Body.Close()

	res, err = ioutil.ReadAll(result.Body)
	return
}

func SendHTTPMethodRequest(method string, urlPath string, body io.Reader, timeOut uint32) (res []byte, err error) {
	httpRequest, err := http.NewRequest(method, urlPath, body)
	if err != nil {
		return
	}

	client := newTimeoutHTTPClient(time.Duration(timeOut) * time.Second)
	result, err := client.Do(httpRequest)
	if err != nil {
		return
	}
	defer result.Body.Close()
	res, err = ioutil.ReadAll(result.Body)
	return
}

func dialHTTPTimeout(timeOut time.Duration) func(net, addr string) (net.Conn, error) {
	return func(network, addr string) (c net.Conn, err error) {
		c, err = net.DialTimeout(network, addr, timeOut)
		if err != nil {
			return
		}
		if timeOut > 0 {
			c.SetDeadline(time.Now().Add(timeOut))
		}
		return
	}
}

func newTimeoutHTTPClient(timeOut time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial: dialHTTPTimeout(timeOut),
		},
	}
}

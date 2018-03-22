package main

import (
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

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

func main() {
	url = "127.0.0.1:9090/test/post"

	client := NewHTTPClient()
	request, err := http.NewRequest("POST", url, strings.NewReader("name=Joe"))
	request.Header.Set("Content-Ttpe", "application/x-www-form-urlencoded")

	if err != nil {
		panic(err)
	}

	response, err := client.Do(request)
	defer response.Body.Close()
	stdout := os.Stdout
	_, err = io.Copy(stdout, response.Body)
	status := response.StatusCode
	body, err := ioutil.ReadAll(response.Body)
	log.Println(string(body))
	log.Println(status)

	http.PostForm
}

func PostForm(url string, data url.Values) (resp *Response, err error) {
	return DefaultClient.PostForm(url, data)
}

// PostForm issues a POST to the specified URL,
// with data's keys and values URL-encoded as the request body.
//
// The Content-Type header is set to application/x-www-form-urlencoded.
// To set other headers, use NewRequest and DefaultClient.Do.
//
// When err is nil, resp always contains a non-nil resp.Body.
// Caller should close resp.Body when done reading from it.
//
// See the Client.Do method documentation for details on how redirects
// are handled.

// PostForm source code

func (c *Client) PostForm(url string, data url.Values) (resp *Response, err error) {
	return c.Post(url, "application/x-www-form-urlencoded", strings.NewReader(data.Encode()))
}

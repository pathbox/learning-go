package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"time"

	log "github.com/segmentio/go-log"
)

type Person struct {
	Name string
	Age  int
}

func main() {
	fmt.Println("NNN")
	url := "http://127.0.0.1:8001"
	m := make(map[string]string)
	m["name"] = "Joe"
	m["age"] = "27"
	PostJson(url, m)
}

func NewHTTPClient() *http.Client {
	transport := &http.Transport{
		Dial: (&net.Dialer{
			Timeout: 10 * time.Second,
		}).Dial,
		TLSHandshakeTimeout:   10 * time.Second,
		ResponseHeaderTimeout: 10 * time.Second,
	}

	client := &http.Client{
		Timeout:   10 * time.Second,
		Transport: transport,
	}
	return client
}

func PostJson(url string, paramMap map[string]string) (*http.Response, error) {
	// var r http.Request
	// r.ParseForm()

	// for k, v := range paramMap {
	// 	r.Form.Add(k, v)
	// }

	// payload := strings.TrimSpace(r.Form.Encode())
	// request, err := http.NewRequest("POST", url, strings.NewReader(payload))

	reqBody, _ := json.Marshal(paramMap)
	request, err := http.NewRequest("POST", url, bytes.NewReader(reqBody))
	if err != nil {
		log.Error("PostJson", err)
		return nil, err
	}

	client := NewHTTPClient()
	request.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(request)
	if err != nil {
		log.Error("PostJson", err)
	}

	return resp, err
}

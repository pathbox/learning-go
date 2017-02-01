package main

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/http/httptrace"
	"os"
)

func main() {
	server := httptest.NewServer(http.HandlerFunc(http.NotFound))
	defer server.Close()
	c := http.Client{}
	req, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		panic(err)
	}

	trace := &httptrace.ClientTrace{
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Println("Got Conn")
		},
		ConnectStart: func(network, addr string) {
			fmt.Println("Dial start")
		},
		ConnectDone: func(network, addr string, err error) {
			fmt.Println("Dial done")
		},
		GotFirstResponseByte: func() {
			fmt.Println("First response byte!")
		},
		WroteHeaders: func() {
			fmt.Println("Wrote headers")
		},
		WroteRequest: func(wr httptrace.WroteRequestInfo) {
			fmt.Println("Wrote request", wr)
		},
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))
	fmt.Println("Starting request!")
	resp, err := c.Do(req)
	if err != nil {
		panic(err)
	}
	io.Copy(os.Stdout, resp.Body)
	fmt.Println("Done!")
}

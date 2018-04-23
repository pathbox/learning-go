package main

import (
	"crypto/tls"
	"io/ioutil"
	"net"
	"net/http"
	"testing"
)

var httpServer *http.Server
var httpsServer *http.Server

type testHandler struct {
}

func (th *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("OK.\n"))
}

func startHTTPServer() {
	if httpServer != nil {
		return
	}

	httpServer = &http.Server{
		Handler: &testHandler{},
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		panic(err)
	}

	go func() {
		err := httpServer.Serve(listener)
		if err != nil {
			panic(err)
		}
	}()
}

func startHTTPSServer() {
	if httpsServer != nil {
		return
	}

	httpsServer = &http.Server{
		Handler: &testHandler{},
	}

	listener, err := net.Listen("tcp", ":8443")
	if err != nil {
		panic(err)
	}

	go func() {
		err := httpServer.ServeTLS(listener, "server.crt", "server.key")
		if err != nil {
			panic(err)
		}
	}()
}

func sendRequest(client *http.Client, addr string) {
	res, err := client.Get(addr)
	if err != nil {
		panic(err)
	}

	if res.StatusCode != 200 {
		panic("request failed")
	}

	_, err = ioutil.ReadAll(res.Body)
	if err != nil {
		panic(err)
	}

	err = res.Body.Close()
	if err != nil {
		panic(err)
	}
}

func BenchmarkHTTP(b *testing.B) {
	startHTTPServer()

	client := &http.Client{}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sendRequest(client, "http://127.0.0.1:8080/")
	}
}

func BenchmarkHTTPNoKeepAlive(b *testing.B) {
	startHTTPServer()

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sendRequest(client, "http://127.0.0.1:8080/")
	}
}

func BenchmarkHTTPSNoKeepAlive(b *testing.B) {
	startHTTPSServer()

	client := &http.Client{
		Transport: &http.Transport{
			DisableKeepAlives: true,
			TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		},
	}

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		sendRequest(client, "https://127.0.0.1:8443/")
	}
}

/*
BenchmarkHTTP-4              10,000    189,912 ns/op    5,736 B/op   70 allocs/op
BenchmarkHTTPNoKeepAlive-4    5,000    359,027 ns/op   17,204 B/op  123 allocs/op
BenchmarkHTTPSNoKeepAlive-4     300  4,052,008 ns/op  116,289 B/op  843 allocs/op
*/

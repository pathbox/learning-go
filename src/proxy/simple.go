package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

func main() {
	targetUrl, err := url.Parse("http://httpbin.org")
	if err != nil {
		panic("bad url")
	}

	proxy := httputil.NewSingleHostReverseProxy(targetUrl)

	http.Handle("/", proxy)
	log.Println("Start serving on port 1234")

	http.ListenAndServe(":1234", nil)
	os.Exit(0)
}

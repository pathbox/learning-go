package main

import (
	"net/http"
	"time"
)

const (
	timeout    = time.Duration(1 * time.Second)
	timeoutMsg = "Your request has timed out"
)

func myTimeoutHandler(h http.Handler) http.Handler {
	return http.TimeoutHandler(h, timeout, timeoutMsg)
}

func main() {
	indexHandler := http.HandlerFunc(index)

	http.Handle("/", myTimeoutHandler(indexHandler))
	http.ListenAndServe(":9090", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

package main

import (
	"fmt"
	"html"
	"net/http"

	"golang.org/x/net/http2"
	"log"
)

func main() {
	var server http.Server

	http2.VerboseLogs = true
	server.Addr = "127.0.0.1:9090"

	http2.ConfigureServer(&server, &http2.Server{})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "URL: %q\n", html.EscapeString(r.URL.Path))
		ShowRequestInfoHandler(w, r)
	})

	server.ListenAndServe() //不启用 https 则默认只支持http1.x
	//log.Fatal(server.ListenAndServeTLS("localhost.cert", "localhost.key"))
}

func ShowRequestInfoHandler(w http.ResponseWriter, r *http.Request) {

	log.Println("here")
	//    return

	w.Header().Set("Content-Type", "text/plain")

	fmt.Fprintf(w, "Method: %s\n", r.Method)
	fmt.Fprintf(w, "Protocol: %s\n", r.Proto)
	fmt.Fprintf(w, "Host: %s\n", r.Host)
	fmt.Fprintf(w, "RemoteAddr: %s\n", r.RemoteAddr)
	fmt.Fprintf(w, "RequestURI: %q\n", r.RequestURI)
	fmt.Fprintf(w, "URL: %#v\n", r.URL)
	fmt.Fprintf(w, "Body.ContentLength: %d (-1 means unknown)\n", r.ContentLength)
	fmt.Fprintf(w, "Close: %v (relevant for HTTP/1 only)\n", r.Close)
	fmt.Fprintf(w, "TLS: %#v\n", r.TLS)
	fmt.Fprintf(w, "\nHeaders:\n")

	r.Header.Write(w)
}

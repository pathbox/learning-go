package main

import (
	"flag"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

func main() {
	var src, dst string
	flag.Parse()
	args := flag.Args()
	if len(args) >= 1 {
		dst = args[0]
	} else {
		dst = "127.0.0.1:8080"
	}

	if len(args) == 2 {
		src = args[1]
	} else {
		src = ":8888"
	}
	u, e := url.Parse(dst)
	if e != nil {
		log.Fatal("Bad destination.")
	}
	h := httputil.NewSingleHostReverseProxy(u)
	s := &http.Server{
		Addr:           src,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Fatal(s.ListenAndServe())
}

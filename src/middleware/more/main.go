package main

import (
	"bytes"
	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
	"net/http"
	"os"
)

func enforceXMLHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(400), 400)
			return
		}
		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		if http.DetectContentType(buf.Bytes()) != "text/html; charset=utf-8" {
			http.Error(w, http.StatusText(415), 415)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func myLoggingHandler(h http.Handler) http.Handler {
	logFile, err := os.OpenFile("server.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return handlers.LoggingHandler(logFile, h)
}

func main() {
	indexHandler := http.HandlerFunc(index)
	authHandler := httpauth.SimpleBasicAuth("username", "password")
	http.Handle("/", myLoggingHandler(authHandler(enforceXMLHandler(indexHandler))))
	http.ListenAndServe(":9090", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

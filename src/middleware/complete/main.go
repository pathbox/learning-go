package main

import (
	"bytes"
	"net/http"
	"os"

	"github.com/goji/httpauth"
	"github.com/gorilla/handlers"
)

func enforeceXMLHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength == 0 {
			http.Error(w, http.StatusText(400), 400)
			return
		}

		buf := new(bytes.Buffer)
		buf.ReadFrom(r.Body)
		if http.DetectContentType(buf.Bytes()) != "text/xml; charset=utf-8" {
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
	return handlers.LoggingHandler(logFiler, h)
}

func main() {
	indexHandler := http.HandlerFunc(index)
	authHandler := httpauth.SimpleBasicAuth("admin", "password")

	http.Handle("/", myLoggingHandler(authHandler(enforeceXMLHandler(indexHandler))))
	http.ListenAndServe("9090", nil)
}

func index(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

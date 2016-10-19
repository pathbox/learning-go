package main

import (
	"net/http"

	"github.com/goji/httpauth"
)

func main() {
	finalHandler := http.HandlerFunc(final)
	authHandler := httpauth.SimpleBasicAuth("username", "password")

	http.Handle("/", authHandler(finalHandler)) // middleware is that nesting handlers stack
	http.ListenAndServe(":3000", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World"))
}

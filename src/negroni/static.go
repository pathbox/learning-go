package main

import (
	"fmt"
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	// Example of using a http.FileServer if you want "server-like" rather than "middleware" behavior
	// mux.Handle("/public", http.FileServer(http.Dir("/home/public")))

	n := negroni.New()
	n.Use(negroni.NewStatic(http.Dir("/tmp")))
	n.UseHandler(mux)

	http.ListenAndServe(":3002", n)
}

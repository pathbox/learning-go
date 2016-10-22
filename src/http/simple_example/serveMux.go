package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
  h := http.NewServeMux()

  // http.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
  //   fmt.Fprintf(w, "Hello, you hit foo")
  // })

  h.HandleFunc("/foo", func(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintf(w, "Hello, you hit foo")
  })

  h.HandleFunc("/bar", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, you hit bar!")
	})

	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		fmt.Fprintln(w, "You're lost, go home")
	})

  log.Fatal(http.ListenAndServe(":9090", h))

}

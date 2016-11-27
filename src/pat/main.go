package main

import (
	"io"
	"log"
	"net/http"

	"github.com/bmizerany/pat"
)

func HelloServer(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "hello, "+req.URL.Query().Get(":name")+"!\n")
}

func main() {
	m := pat.New()
	m.Get("/hello/:name", http.HandlerFunc(HelloServer))

	http.Handle("/", m)
	log.Fatal(http.ListenAndServe(":9090", nil))
}

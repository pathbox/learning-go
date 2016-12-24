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

	n := negroni.New(negroni.NewLogger(), negroni.NewRecovery(), negroni.NewStatic(http.Dir("/tmp")))
	n.UseHandler(mux)

	http.ListenAndServe(":3004", n)
}

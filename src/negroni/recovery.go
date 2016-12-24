package main

import (
	"net/http"

	"github.com/urfave/negroni"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		panic("oh no no no no no")
	})

	n := negroni.New()
	n.Use(negroni.NewRecovery())
	n.UseHandler(mux)

	http.ListenAndServe(":3003", n)
}

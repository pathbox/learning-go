package main

import (
	"net/http"

	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/testutils"
)

func main() {
	// Forwards incoming requests to whatever location URL points to, adds proper forwarding headers
	fwd, _ := forward.New()
	redirect := http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		req.URL = testutils.ParseURI("http://127.0.0.1:63450")
		fwd.ServeHTTP(w, req)
	})

	s := &http.Server{
		Addr:    ":8080",
		Handler: redirect, // request ":8080", in fact request http://127.0.0.1:63450
	}
	s.ListenAndServe()
}

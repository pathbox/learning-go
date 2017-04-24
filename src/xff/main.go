package main

import (
	"net/http"

	"github.com/sebest/xff"
)

func main() {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hello from " + r.RemoteAddr + "\n"))
	})
	xffmw, _ := xff.Default()
	http.ListenAndServe(":9090", xffmw.Handler(handler))
}

// curl -D - -H 'X-Forwarded-For: 42.42.42.42' http://localhost:9090/

package main

import (
	"net/http"

	"github.com/vulcand/oxy/buffer"
	"github.com/vulcand/oxy/forward"
	"github.com/vulcand/oxy/roundrobin"
)

func main() {
	fwd, _ := forward.New()
	lb, _ := roundrobin.New(fwd)

	// buffer will read the request body and will replay the request again in case if forward returned status
	// corresponding to nework error (e.g. Gateway Timeout)
	buffer, _ := buffer.New(lb, buffer.Retry(`IsNetworkError() && Attempts() < 2`))

	lb.UpsertServer(url1)
	lb.UpsertServer(url2)

	s := &http.Server{
		Addr:    ":8080",
		Handler: lb,
	}
	s.ListenAndServe()
}

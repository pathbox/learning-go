package main

import (
	"context"
	"net/http"
)

type ContextAdapter struct {
	ctx     context.Context
	handler ContextHandler
}

func (ca *ContextAdapter) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	ca.handler.ServeHTTPContext(ca.ctx, rw, req)
}

func main() {
	h := &ContextAdapter{
		ctx:     context.Background(),
		handler: middleware(ContextHandlerFunc(handler)),
	}
	http.ListenAndServe(":8080", h)
}

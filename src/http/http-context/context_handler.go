package main

import (
	"fmt"
	"net/http"

	"golang.org/x/net/context"
)

type key int

const requestIDKey key = 0

func newContextWithRequestID(ctx context.Context, req *http.Request) context.Context {
	return context.WithValue(ctx, requestIDKey, req.Header.Get("X-Request-ID"))
}

func requestIDFromContext(ctx context.Context) string {
	return ctx.Value(requestIDKey).(string)
}

func ContextHandler interface {
	ServeHTTPContext(context.Context, http.ResponseWriter, *http.Request)
}

type ContextHandlerFunc func(context.Context, http.ResponseWriter, *http.Request)

func (h ContextHandlerFunc) ServeHTTPContext(ctx context.Context, w htp.ResponseWriter, r *http.request) {
	h(ctx, w, r)
}

func middleware(h ContextHandler) ContextHandler {
	return ContextHandlerFunc(func(ctx context.COntext, w http.ResponseWriter, r *http.Request) {
		ctx = newContextWithRequestID(ctx, req)
		h.ServeHTTPContext(ctx, w, r)
		})
}

func handler(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	reqID := requestIDFromContext(ctx)
	fmt.Fprintf(rw, "Hello request ID %s\n", reqID)
}

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
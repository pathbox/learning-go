var cmap = map[*http.Request]*context.Context{}
var cmapLock sync.Mutex

// Note that we are returning a pointer to the context, not the
// context itself.
func contextFromRequest(req *http.Request) *context.Context {
    cmapLock.Lock()
    defer cmapLock.Unlock()
    return cmap[req]
}

// Necessary wrapper around all handlers.  Must be the first middleware.
func contextHandler(ctx context.Context, h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctx2 := ctx // make a copy of the root context reference
        cmapLock.Lock()
        cmap[req] = &ctx2
        cmapLock.Unlock()

        h.ServeHTTP(rw, req)

        cmapLock.Lock()
        delete(cmap, req)
        cmapLock.Unlock()
    })
}

func middleware(h http.Handler) http.Handler {
    return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
        ctxp := contextFromRequest(req)
        *ctxp = newContextWithRequestID(*ctxp, req)

        h.ServeHTTP(rw, req)
    })
}

func handler(rw http.ResponseWriter, req *http.Request) {
    ctxp := contextFromRequest(req)

    reqID := requestIDFromContext(*ctxp)
    fmt.Fprintf(rw, "Hello request ID %s\n", reqID)
}

func main() {
    h := contextHandler(context.Background(), middleware(http.HandlerFunc(handler)))
    http.ListenAndServe(":8080", h)
}
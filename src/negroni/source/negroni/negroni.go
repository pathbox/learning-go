package negroni

import (
  "log"
  "net/http"
  "os"
)

// Handler handler is an interface that objects can implement to be registered to serve as middleware
// in the Negroni middleware stack.
// ServeHTTP should yield to the next middleware in the chain by invoking the next http.HandlerFunc
// passed in.
//
// If the Handler writes to the ResponseWriter, the next http.HandlerFunc should not be invoked.

type Handler interface {
  ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

type HandlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)

func (h HandlerFunc) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc){
  h(rw, r, next)
}

type middleware struct{
  handler Handler
  next *middleware
}

func(m middleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
  m.handler.ServeHTTP(rw, r, m.next.ServeHTTP)
}

func Wrap(handler http.Handler) Handler {
  return HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next.http.HandlerFunc){
    handler.ServeHTTP(rw, r)
    next(rw, r)
    })
}

type Negroni struct {
  middleware middleware
  handlers []Handler
}

func New(handlers ...Handler) *Negroni{
  return &Negroni{
    handlers: handlers,
    middleware: build(handlers),
  }
}

func Classic() *Negroni{
  return New(NewRecovery(), NewLogger(), NewStatic(http.Dir("public")))
}

func (n *Negroni) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
  n.middleware.ServeHTTP(NewResponseWrite(rw),r)
}

func (n *Negroni) Use(handler Handler) {
  if handler == nil {
    panic("handler cannot be nil")
  }

  n.handlers = append(n.handlers, handler)
  n.middleware = build(n.handlers)
}

func (n *Negroni) UseFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc)){
  n.Use(HandlerFunc(handlerFunc))
}

// UseHandler adds a http.Handler onto the middleware stack. Handlers are invoked in the order they are added to a Negroni.
func (n *Negroni) UseHandler(handler http.Handler) {
  n.Use(Wrap(handler))
}

func (n *Negroni) UseHandlerFunc(handlerFunc func(rw http.ResponseWriter, r *http.Request)) {
  n.UseHandler(http.HandlerFunc(handlerFunc))
}

func (n *Negroni) Run(addr string) {
  l := log.New(os.Stdout, "[negroni] ", 0)
  l.Printf("listening on %s", addr)
  l.Fatal(http.ListenAndServe(addr, n))
}

// Returns a list of all the handlers in the current Negroni middleware chain.
func (n *Negroni) Handlers() []Handler {
  return n.handlers
}

func build(handlers []Handler) middleware {
  var next middleware

  if len(handlers) == 0 {
    return voidMiddleware()
  } else if len(handlers) > 1 {
    next = build(handlers[1:])
  } else {
    next = voidMiddleware()
  }

  return middleware{handlers[0], &next}
}

func voidMiddleware() middleware {
  return middleware{
    HandlerFunc(func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {}),
    &middleware{},
  }
}
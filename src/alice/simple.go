package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/alice"
)

type numberDumper int

func (n numberDumper) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
  fmt.Fprintf(w, "Here's your number: %d\n", n)
}

func logger(h http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
    log.Printf("%s requested %s", r.RemoteAddr, r.URL)
    h.ServeHTTP(w, r)
  })
}

type headerSetter struct {
  key, val string
  handler http.Handler
}

func (hs headerSetter) ServeHTTP(w http.ResponseWriter, r *http.Request)  {
  w.Header().Set(hs.key, hs.val)
  hs.handler.ServeHTTP(w, r)
}

func newHeaderSetter(key, val string) func(http.Handler) http.Handler {
  return func(h http.Handler) http.Handler{
    return headerSetter{key, val, h}
  }
}
func main() {
  h := http.NewServeMux()

  h.Handle("/one", numberDumper(1))
  h.Handle("/two", numberDumper(2))
	h.Handle("/three", numberDumper(3))
	h.Handle("/four", numberDumper(4))

  fiveHS := newHeaderSetter("X-FIVE", "the best number")
  h.Handle("/five", fiveHS(numberDumper(5)))

  h.HandleFunc("/", func (w http.ResponseWriter, r *http.Request)  {
    w.WriteHeader(404)
    fmt.Fprintln(w, "That's not a supported number!")
  })

  chain := alice.New(
    newHeaderSetter("X-FOO", "BAR"),
    newHeaderSetter("X-BAR", "BUZ"),
    logger,
  ).Then(h)

  err := http.ListenAndServe(":9999", chain)
  log.Fatal(err)
}

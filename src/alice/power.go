package main

import (
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"gopkg.in/throttled/throttled.v1"
)

func timeoutHandler(h http.Handler) http.Handler {
  return http.TimeoutHandler(h, 1*time.Second, "timed out")
}

func myApp(w http.ResponseWriter, r *http.Request) {
  w.Write([]byte("Hello World!"))
}

func main() {
  th := throttled.Interval(throttled.PerSec(10), 1, &throttled.VaryBy{Path: true}, 50)
  myHandler := http.HandlerFucn(myApp)

  chain := alice.New(th.Throttle, timeoutHandler, nosurf.NewPure).Then(myHandler)
  http.ListenAndServe(":8000", chain)
}

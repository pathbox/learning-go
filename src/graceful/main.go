package main

import (
	"fmt"
	"github.com/codegangsta/negroni"
	"gopkg.in/tylerb/graceful.v1"
	"net/http"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		fmt.Fprintf(w, "Welcome to the home page!")
	})

	n := negroni.Classic()
	n.UseHandler(mux)
	graceful.Run(":9090", 10*time.Second, n)
}

// Behaviour

// When Graceful is sent a SIGINT or SIGTERM (possibly from ^C or a kill command), it:

// Disables keepalive connections.
// Closes the listening socket, allowing another process to listen on that port immediately.
// Starts a timer of timeout duration to give active requests a chance to finish.
// When timeout expires, closes all active connections.
// Closes the stopChan, waking up any blocking goroutines.
// Returns from the function, allowing the server to terminate

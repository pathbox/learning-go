package main

import (
	"log"
	"net/http"
	_ "net/http/pprof"
)

func main() {
	http.HandleFunc("/", sayHello)
	// http.Handle("debug/pprof", http.HandlerFunc(Index))
	log.Println(http.ListenAndServe(":6060", nil))
}

func sayHello(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello World!"))
}

// go tool pprof http://localhost:6060/debug/pprof/heap
// Or to look at a 30-second CPU profile:

// go tool pprof http://localhost:6060/debug/pprof/profile
// Or to look at the goroutine blocking profile, after calling runtime.SetBlockProfileRate in your program:

// go tool pprof http://localhost:6060/debug/pprof/block

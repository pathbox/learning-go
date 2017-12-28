package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/tidwall/evio"
)

func main() {
	var port int
	var udp bool
	flag.IntVar(&port, "port", 9600, "server port")
	flag.BoolVar(&udp, "udp", false, "listen on udp")
	flag.Parse()

	var events evio.Events
	events.Serving = func(srv evio.Server) (action evio.Action) {
		log.Printf("echo server started on port %d", port)
		return
	}
	events.Opened = func(id int, info evio.Info) (out []byte, opts evio.Options, action evio.Action) {
		log.Printf("opened: %d: %+v", id, info)
		return
	}
	events.Closed = func(id int, err error) (action evio.Action) {
		log.Printf("closed: %d", id)
		return
	}
	events.Data = func(id int, in []byte) (out []byte, action evio.Action) {
		out = in
		return
	}
	scheme := "tcp"
	if udp {
		scheme = "udp"
	}
	log.Fatal(evio.Serve(events, fmt.Sprintf("%s://:%d", scheme, port)))
}

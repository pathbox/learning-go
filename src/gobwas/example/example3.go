package main

import (
	"log"
	"net"

	ws "github.com/gobwas/ws"
)

func main() {
	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatal(err)
	}

	u := ws.Upgrader{
		OnHeader: func(key, value []byte) (err error, code int) {
			log.Printf("non-websocket header: %q=%q", key, value)
			return
		},
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
		}

		_, err = u.Upgrade(conn)
		if err != nil {
			// handle error
		}
	}
}

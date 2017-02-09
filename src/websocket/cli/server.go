package main

import (
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

func echoHadnler(ws *websocket.Conn) {
	io.Copy(ws, ws)
}

func main() {
	http.Handle("/echo", websocket.Handler(echoHadnler))
	err := http.ListenAndServe(":9090", nil)
	if err != nil {
		panic(err)
	}
}

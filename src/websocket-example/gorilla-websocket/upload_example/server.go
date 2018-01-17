package main

import (
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{}

func main() {

	addr := "127.0.0.1:9090"

	http.HandleFunc("/upload", uploadHandler)
	log.Println("Server start...")
	log.Fatal(http.ListenAndServe(addr, nil))
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("uploadHandler")

	path := "./receive.csv"
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				log.Panic(err)
			}
		} else {
			file.Write(message)
		}

	}
}

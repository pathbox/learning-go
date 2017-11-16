package main

import (
	"flag"
	"fmt"
	"net/url"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:12345", "http service address")

func main() {
	u := url.URL{Scheme: "ws", Host: *addr, Path: "/ws"}

	size := 100

	for i := 0; i < size; i++ { // 起100个goroutine进行 websocket连接,实际并发的时候会达到大约200个tcp连接的占用,这个还不清楚原因
		go func() {
			for {
				multiClient(u)
			}

		}()
	}
	time.Sleep(time.Second * 1000) // 简单的让main goroutine 阻塞

}

func multiClient(url url.URL) {
	fmt.Println("Start a new client")
	var dialer *websocket.Dialer
	conn, res, err := dialer.Dial(url.String(), nil)
	fmt.Println("response", res)
	if err != nil {
		fmt.Println(err)
		return
	}
	go timeWriter(conn)

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("read:", err)
			return
		}

		fmt.Printf("received: %s\n", message)
	}
}

func timeWriter(conn *websocket.Conn) {
	for {
		time.Sleep(time.Second * 1)
		conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
	}
}

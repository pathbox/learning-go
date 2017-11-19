package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "59.110.127.112:9090", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// interrupt 是一个chan, 并且处于阻塞状态, 当系统接收到中断信号后,会传入这个chan,然后就会执行中断的处理逻辑

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	size := 100
	for {
		select {
		case <-interrupt:
			os.Exit(-1)
		default:
			for i := 0; i < size; i++ {
				go func() {
					multiClient(u)
				}()
			}
			time.Sleep(time.Second * 1)
		}
	}

	// time.Sleep(time.Second * 60) // 简单的让main goroutine 阻塞

}

func multiClient(url url.URL) {

	var dialer *websocket.Dialer
	conn, _, err := dialer.Dial(url.String(), nil)
	defer conn.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	time.Sleep(time.Second * 1)
	// timeWriter(conn)

	// _, message, err := conn.ReadMessage()
	// if err != nil {
	// 	fmt.Println("client read:", err)
	// 	return
	// }
	// time.Sleep(time.Second * 1)
	// fmt.Printf("received: %s\n", message)
	// time.Sleep(time.Second * 1)

}

func timeWriter(conn *websocket.Conn) {
	conn.WriteMessage(websocket.TextMessage, []byte(time.Now().Format("2006-01-02 15:04:05")))
}

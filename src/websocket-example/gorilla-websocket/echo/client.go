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

var addr = flag.String("addr", "localhost:9090", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// interrupt 是一个chan, 并且处于阻塞状态, 当系统接收到中断信号后,会传入这个chan,然后就会执行中断的处理逻辑

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/echo"}
	log.Printf("connecting to %s", u.String())

	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()
	done := make(chan struct{})

	go func() {
		defer c.Close()
		defer close(done)
		for { // 循环处理 收到的数据
			_, message, err := c.ReadMessage() // 从conn中读取数据
			if err != nil {
				log.Println("read: ", err)
				return
			}
			log.Printf("recv: %s", message)
		}
	}()

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case t := <-ticker.C: // 每隔一秒发送数据给服务端
			// err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			fmt.Println(t) // t 就是类似 2017-11-16 11:37:30.305432349 +0800 CST 这样的时间数据
			err := c.WriteMessage(websocket.TextMessage, []byte("Echo~ 你好"))
			if err != nil {
				log.Println("write:", err)
				return
			}
		case <-interrupt:
			log.Println("interrupt")
			// To cleanly close a connection, a client should send a close
			// frame and wait for the server to close the connection.
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			select {
			case <-done:
				// 最后回调做一些事情
			case <-time.After(time.Second):
				// 超时一秒
			}
			c.Close()
			return
		}
	}
}

package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
)

// var addr = flag.String("addr", "59.110.127.112:9090", "http service address")
var addr = flag.String("addr", "wsecho.com:9090", "http service address")

func main() {
	flag.Parse()
	log.SetFlags(0)

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	// interrupt 是一个chan, 并且处于阻塞状态, 当系统接收到中断信号后,会传入这个chan,然后就会执行中断的处理逻辑

	u := url.URL{Scheme: "wss", Host: *addr, Path: "/echo"}
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

}

func multiClient(url url.URL) {
	pool := x509.NewCertPool()
	caCertPath := "./ca_key/ca.crt"

	caCrt, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		fmt.Println("ReadFile err:", err)
		return
	}

	pool.AppendCertsFromPEM(caCrt) // 客户端添加ca证书

	dialer := &websocket.Dialer{
		TLSClientConfig: &tls.Config{RootCAs: pool},
	}

	conn, _, err := dialer.Dial(url.String(), nil)
	defer conn.Close()
	if err != nil {
		fmt.Println("Conn: ", err)
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

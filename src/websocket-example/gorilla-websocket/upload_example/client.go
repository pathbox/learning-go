package main

import (
	"flag"
	"io"
	"log"
	"net/url"
	"os"

	"github.com/gorilla/websocket"
)

var addr = flag.String("addr", "localhost:9090", "http service address")

func main() {

	from := "/home/user/example.csv"
	file, err := os.Open(from)
	if err != nil {
		log.Panic(err)
	}
	flag.Parse()

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/upload"}
	log.Printf("connecting to %s", u.String())

	client, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	// defer client.Close()

	buf := make([]byte, 1024)

	// 循环读取文件,直到结尾
	for {
		n, err := file.Read(buf)
		if err != nil {
			if err == io.EOF {
			} else {
				log.Panic(err)
			}
			client.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			break
		} else {
			client.WriteMessage(websocket.TextMessage, buf[:n]) // 将每次读取到的字节传到websocket另一端
		}

	}

	log.Println("websocket upload file success")

}

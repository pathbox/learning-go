package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/azhao1981/go-socket.io"
)

func main() {
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}

	server.On("connection", func(so socketio.Socket) { // 前端进入界面就会创建connection
		log.Println("on connection")
		so.Join("chat")
		so.On("chat message", func(msg string) {
			m := make(map[string]interface{})
			m["a"] = "hello"
			e := so.Emit("morning", m)

			fmt.Println("\n\n")

			b := make(map[string]string)
			b["u-a"] = "中文内容" //这个不能是中文
			m["b-c"] = b
			e = so.Emit("golang", m)
			log.Println(e)

			log.Println("emit error:", so.Emit("chat message", msg+"back")) // 给 chat message 名称的渠道socket发msg
			so.BroadcastTo("chat", "chat message", msg)
		})

		so.On("chat message with ack", func(msg string) string {
			return msg
		})

		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error: ", err)
	})

	http.Handle("/socket.io/", server)                     // socket.io 服务
	http.Handle("/", http.FileServer(http.Dir("./asset"))) // 网站服务
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil)) // listen 5000
}

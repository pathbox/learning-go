package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/googollee/go-engine.io.v1"
)

func main() {
	server, _ := engineio.NewServer(nil)

	go func() {
		for {
			conn, _ := server.Accept() // Accept成功了，才会继续往下走，建立一个goroutine，处理这个连接，要不会一直阻塞在Accept，等待新的连接到来。所以，Accept需要写在死循环中，要不只会生成一个连接，逻辑就结束了
			go func() {
				defer conn.Close()
				for {
					t, r, _ := conn.NextReader()
					b, _ := ioutil.ReadAll(r)
					r.Close()

					w, _ := conn.NextWriter(t)
					w.Write(b)
					w.Close()
				}
			}()
		}

	}()

	http.Handle("/engine.io/", server)
	log.Println("Serving at localhost:5000...")
	log.Fatal(http.ListenAndServe(":5000", nil))
}

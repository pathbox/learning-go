package main

import (
	"log"
	"net/http"
)

func main() {
	// 在 8000 端口启动服务器
	// 确切地说，如何运行HTTP/1.1服务器。

	srv := &http.Server{Addr: ":8000", Handler: http.HandlerFunc(handle)}
	// 用TLS启动服务器，因为我们运行的是http/2，它必须是与TLS一起运行。
	// 确切地说，如何使用TLS连接运行HTTP/1.1服务器。
	log.Printf("Serving on https://0.0.0.0:8000")
	log.Fatal(srv.ListenAndServeTLS("server.crt", "server.key"))
}

func handle(w http.ResponseWriter, r *http.Request) {
	// 记录请求协议
	log.Printf("Got connection: %s", r.Proto)
	// 向客户发送一条消息
	w.Write([]byte("Hello"))
}

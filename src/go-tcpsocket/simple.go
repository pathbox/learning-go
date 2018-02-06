package main

import (
	"log"
	"net"
)

func main() {
	go ServerStart()
	select {}
}

func ServerStart() {
	ln, err := net.Listen("tcp", ":9090") // 监听建立连接
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := ln.Accept() // 循环 Accpet(), 阻塞直到有新的连接建立，获取该conn
		if err != nil {
			log.Println(err)
			continue
		}
		go handleConnection(conn) // 处理conn 连接
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	// do something
}

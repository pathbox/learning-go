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
		conn, err := ln.Accept() // 循环 Accpet()
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

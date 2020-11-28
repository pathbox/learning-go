package main

import (
	"fmt"
	"net"
)

func main() {
	listen, err := net.Listen("tcp", ":9090")
	if err != nil {
		panic(err)
	}
	fmt.Println("Listening")
	for {
		conn, err := listen.Accept()
		if err != nil {
			panic(err)
		}
		go handleConn(conn) // 模式 goroutine-per-connection
	}
}
func handleConn(conn net.Conn) {
	defer conn.Close()
	packet := make([]byte, 1024)
	for {
		n, err := conn.Read(packet)
		if err != nil {
			panic(err)
		}
		_, _ = conn.Write(packet[:n])
	}
}

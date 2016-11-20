package main

import (
	"fmt"
	"net"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		// read from the connection

		// write to the connection
	}
}

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		fmt.Println("Listen error:", err)
		return
	}
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("Accept error:", err)
			break
		}
		go handleConn(c)
	}
}

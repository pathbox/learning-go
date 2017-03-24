//粘包问题演示客户端
package main

import (
	"./protocol"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

func sender(conn net.Conn) {
	for i := 0; i < 100; i++ {
		words := "{\"Id\": " + strconv.Itoa(i) + ",\"Name\":\"golang\",\"Message\":\"message\"}"
		conn.Write(protocol.Packet([]byte(words)))
	}
	fmt.Println("send over")
}

func main() {
	server := "127.0.0.1:9090"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	defer conn.Close()

	fmt.Println("connect success")

	go sender(conn)

	// for {
	// 	time.Sleep(1 * 10)
	// }
	time.Sleep(3 * time.Second)
}

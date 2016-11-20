package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":9090"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	fmt.Println("Listenning...")
	checkError(err)
	for {
		conn, err := listener.Accept()
		defer conn.Close() // we're finished with this client
		if err != nil {
			continue
		}
		daytime := time.Now().String()
		conn.Write([]byte(daytime)) // don't care about return value
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

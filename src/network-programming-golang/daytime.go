package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	service := ":9099"

	tcpAddr, err := net.ResolveTCPAddr("tcp", service) // build the tcpAddr for listen
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
	fmt.Println("Start daytime server")
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}

		daytime := time.Now().String()
		conn.Write([]byte(daytime))
		conn.Close() // just one connection close, not the listener
	}
}

// It is tcp server, you can't get data from Browser, just use tcp client or curl

func checkError(err error) {

	if err != nil {

		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())

		os.Exit(1)

	}

}

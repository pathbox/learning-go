package main

import (
	"encoding/asn1"

	"fmt"

	"net"

	"os"

	"time"
)

func main() {

	service := ":9090"

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)

	checkError(err)

	listener, err := net.ListenTCP("tcp", tcpAddr)

	checkError(err)

	for {

		conn, err := listener.Accept()

		if err != nil {

			continue

		}

		daytime := time.Now()

		// Ignore return network errors.

		mdata, _ := asn1.Marshal(daytime)
		fmt.Println(string(mdata))
		conn.Write(mdata)

		conn.Close() // we're finished

	}

}

func checkError(err error) {

	if err != nil {

		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())

		os.Exit(1)

	}

}

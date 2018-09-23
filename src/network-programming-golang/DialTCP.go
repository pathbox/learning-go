package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {

		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])

		os.Exit(1)

	}
	service := os.Args[1]

	tcpAddr, err := net.ResolveTCPAddr("tcp", service)
	if err != nil {
		panic(err)
	}

	fmt.Println("tcpAddr: ", tcpAddr)

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)

	_, err = conn.Write([]byte("HEAD / HTTP/1.1\r\n\r\n")) // you can pass json data to the service
	checkError(err)

	result, err := ioutil.ReadAll(conn)
	checkError(err)

	fmt.Println(string(result))
}

func checkError(err error) {

	if err != nil {

		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())

		os.Exit(1)

	}

}

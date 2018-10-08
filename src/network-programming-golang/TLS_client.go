package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: ", os.Args[0], "host:port")
		os.Exit(1)
	}

	cert, err := tls.LoadX509KeyPair("cary.pathbox.pem", "private.pem")
	checkError(err)
	config := tls.Config{Certificates: []tls.Certificate{cert}}

	now := time.Now()
	config.Time = func() time.Time { return now }
	config.Rand = rand.Reader

	service := os.Args[1]
	fmt.Println("The service address is: ", service)
	conn, err := tls.Dial("tcp", service, &config)
	checkError(err)

	for n := 0; n < 10; n++ {
		fmt.Println("Waiting...")
		conn.Write([]byte("Hello " + string(n+48)))

		var buf [512]byte
		_, err := conn.Read(buf[0:])
		checkError(err)

		fmt.Println(string(buf[:]))
	}
	os.Exit(6)
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

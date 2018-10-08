// Encryption/decryption schemes are of limited use if you have to do all the heavy lifting yourself. The most popular mechanism on the internet to give support for encrypted message passing is currently TLS (Transport Layer Security) which was formerly SSL (Secure Sockets Layer).

// In TLS, a client and a server negotiate identity using X.509 certificates. Once this is complete, a secret key is invented between them, and all encryption/decryption is done using this key. The negotiation is relatively slow, but once complete a faster private key mechanism is used
package main

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"net"
	"os"
	"time"
)

func main() {
	cert, err := tls.LoadX509KeyPair("cary.pathbox.pem", "private.pem")
	checkError(err)
	config := tls.Config{Certificates: []tls.Certificate{cert}}

	now := time.Now()
	config.Time = func() time.Time { return now }
	config.Rand = rand.Reader

	service := "0.0.0.0:9001"

	listener, err := tls.Listen("tcp", service, &config)
	checkError(err)
	fmt.Println("Listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Accepted")
		go handleClient(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close() // handle every conn, then close the conn

	var buf [512]byte
	for {
		fmt.Println("Trying to read")
		n, err := conn.Read(buf[0:])
		if err != nil {
			fmt.Println(err)
		}
		conn.Write([]byte(`I get it, echo!`))
		_, err2 := conn.Write(buf[0:n])
		if err2 != nil {
			return
		}
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Fatal error ", err.Error())
		os.Exit(1)
	}
}

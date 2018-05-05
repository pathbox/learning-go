package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	sock, err := net.Dial("tcp", ":12345")
	if err != nil {
		log.Fatalln("fail to dial ':12345':", err)
	}

	buffer := make([]byte, 64)
	for {
		n, err := sock.Read(buffer)
		if err != nil {
			log.Fatalln("fail to read socket:", err)
		}
		fmt.Printf("[Client] Received %d bytes, '%s'\n", n, string(buffer[:n]))
		_, err = sock.Write([]byte("pong"))
		if err != nil {
			log.Fatalln("fail to write 'pong' to the socket:", err)
		}

		fmt.Printf("[Client] Sent 'ping'\n")
	}

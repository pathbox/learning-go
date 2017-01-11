package main

import (
	"log"
	"net"
	"time"
)

func main() {
	message := "1234567890"
	log.Println("begin dial...")
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		log.Println("dial error: ", err)
		return
	}
	defer conn.Close()
	log.Println("dial ok")

	time.Sleep(time.Second * 2)
	conn.Write([]byte(message))
	time.Sleep(time.Second * 100000)
}

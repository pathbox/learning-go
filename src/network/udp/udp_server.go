package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	service := ":9090"
	updAddr, err := net.ResolveUDPAddr("udp4", service)
	checkError(err)
	conn, err := net.ListenUDP("udp", updAddr)
	checkError(err)
	for {
		handleClient(conn)
	}
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error %v", err.Error())
		os.Exit(1)
	}
}

func handleClient(conn *net.UDPConn) {
	var buf [512]byte
	res, addr, err := conn.ReadFromUDP(buf[0:])
	if err != nil {
		log.Printf("Read from UDP failed, err: %v", err)
		return
	}
	fmt.Fprintf(os.Stderr, "response: %v \n", res)
	fmt.Println(addr)
	daytime := time.Now().String()
	conn.WriteToUDP([]byte(daytime), addr)
}

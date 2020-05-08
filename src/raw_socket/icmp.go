package main

import (
	"fmt"
	"net"

	"golang.org/x/net/icmp"
)

func main() {
	netaddr, _ := net.ResolveIPAddr("ip4", "180.101.49.12")
	conn, err := net.ListenIP("ip4:icmp", netaddr)
	if err != nil {
		fmt.Println("err:", err)
		panic(err)
	}
	for {
		buf := make([]byte, 1024)
		n, addr, _ := conn.ReadFrom(buf)
		msg, _ := icmp.ParseMessage(1, buf[0:n])
		fmt.Println(n, addr, msg.Type, msg.Code, msg.Checksum)
	}
}

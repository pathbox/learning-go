package main

import (
	"fmt"
	"net"

	"./tcp"
)

func main() {
	netaddr, _ := net.ResolveIPAddr("ip4", "0.0.0.0")
	conn, _ := net.ListenIP("ip4:tcp", netaddr)
	for {
		buf := make([]byte, 1480)
		n, addr, _ := conn.ReadFrom(buf)
		tcpheader := tcp.NewTCPHeader(buf[0:n])
		fmt.Println(n, addr, tcpheader)
	}
}

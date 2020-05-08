package main

import (
	"golang.org/x/net/ipv4"
	"net"
	"syscall"
)

func main() {
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	addr := syscall.SockaddrInet4{
		Port: 0,
		Addr: [4]byte{172, 17, 0, 3},
	}
	yipHeader := ipv4.Header{
		Version:  4,
		Len:      20,
		TotalLen: 20, // 20 bytes for IP, 10 for ICMP
		TTL:      64,
		Protocol: 6, // TCP
		Dst:      net.IPv4(172, 17, 0, 3),
		Src:      net.IPv4(172, 17, 0, 99),
	}
	payload, _ := yipHeader.Marshal()
	syscall.Sendto(fd, payload, 0, &addr)
}
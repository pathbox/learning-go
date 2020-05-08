package main

import (
	"github.com/mushroomsir/blog/examples/util"
	"log"
	"net"
	"strconv"
	"strings"
)

func main() {
	local := "127.0.0.1"
	remote := "172.17.0.3"
	conn, _ := net.Dial("ip4:tcp", remote)
	ycpHeader:= util.TCPHeader{
		Source:      17663, // Random ephemeral port
		Destination: 8020,
		SeqNum:      2,
		AckNum:      0,
		DataOffset:  5,      // 4 bits
		Reserved:    0,      // 3 bits
		ECN:         0,      // 3 bits
		Ctrl:        2,      // 6 bits (000010, SYN bit set)
		Window:      0xaaaa, // size of your receive window
		Checksum:    0,      // Kernel will set this if it's 0
		Urgent:      99,
	}
	data := ycpHeader.Marshal()
	ycpHeader.Checksum = util.Csum(data, to4byte(local), to4byte(remote))
	data = ycpHeader.Marshal()
	data=append(data,[]byte("xxx")...)
	conn.Write(data)
}

func to4byte(addr string) [4]byte {
	parts := strings.Split(addr, ".")
	b0, err := strconv.Atoi(parts[0])
	if err != nil {
		log.Fatalf("to4byte: %s (latency works with IPv4 addresses only, but not IPv6!)\n", err)
	}
	b1, _ := strconv.Atoi(parts[1])
	b2, _ := strconv.Atoi(parts[2])
	b3, _ := strconv.Atoi(parts[3])
	return [4]byte{byte(b0), byte(b1), byte(b2), byte(b3)}
}
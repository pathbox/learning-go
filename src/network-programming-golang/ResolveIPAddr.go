package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {

		fmt.Fprintf(os.Stderr, "Usage: %s hostname\n", os.Args[0])

		fmt.Println("Usage: ", os.Args[0], "hostname")

		os.Exit(1)

	}
	name := os.Args[1]

	addr, err := net.ResolveIPAddr("ip", name) // ip ip4 ip6
	if err != nil {
		panic(err)
	}
	fmt.Println("Resolved address is ", addr.String())
}

// go run ResolveIPAddr.go www.baidu.com
// Resolved address is  61.135.169.125

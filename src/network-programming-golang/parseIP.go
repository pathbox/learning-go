package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s ip-addr\n", os.Args[0])
		os.Exit(1)
	}
	name := os.Args[1]
	addr := net.ParseIP(name)
	if addr == nil {

		fmt.Println("Invalid address")

	} else {

		fmt.Println("The address is ", addr.String())

	}

	os.Exit(0)
}

// IP 127.0.0.1
// IP 0:0:0:0:0:0:0:1
// go run parseIP.go 20:0:3:0:0:e0:0:1
// The address is  20:0:3::e0:0:1

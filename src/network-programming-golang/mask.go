package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s dotted-ip-addr\n", os.Args[0])

		os.Exit(1)
	}

	dotAddr := os.Args[1]

	addr := net.ParseIP(dotAddr)
	if addr == nil {
		fmt.Println("Invalid address")
		os.Exit(1)
	}

	mask := addr.DefaultMask()
	network := addr.Mask(mask)
	ones, bits := mask.Size()
	fmt.Println("Address is ", addr.String(),

		" Default mask length is ", bits,

		"Leading ones count is ", ones,

		"Mask is (hex) ", mask.String(),

		" Network is ", network.String())

}

// go run mask.go 192.168.1.1

// Address is  192.168.1.1  Default mask length is  32 Leading ones count is  24 Mask is (hex)  ffffff00  Network is  192.168.1.0

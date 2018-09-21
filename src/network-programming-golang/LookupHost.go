package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s hostname\n", os.Args[0])
		os.Exit(1)
	}

	name := os.Args[1]

	addrs, err := net.LookupHost(name)
	if err != nil {

		fmt.Println("Error: ", err.Error())

		os.Exit(2)

	}
	for _, s := range addrs {
		fmt.Println(s)
	}

}

// go run LookupHost.go www.baidu.com
// 119.75.216.20
// 119.75.213.61

// get multi IP from the hostname

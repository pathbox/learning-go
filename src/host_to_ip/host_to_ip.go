package main

import (
	"fmt"
	"net"
)

func main() {
	addr, err := net.LookupHost("ucloud.cn")
	if err, ok := err.(*net.DNSError); ok {
		if err.Timeout() {
			fmt.Println("operation timed out")
		} else if err.Temporary() {
			fmt.Println("temporary error")
		} else {
			fmt.Println("generic error: ", err)
		}
		return
	}
	fmt.Println(addr)
}

package main

import (
	"fmt"
	"os"
	"syscall"

	"golang.org/x/net/ipv4"
)

func main() {
	fd, _ := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_ICMP)
	f := os.NewFile(uintptr(fd), fmt.Sprintf("fd %d", fd))
	for {
		buf := make([]byte, 1024)
		numRead, err := f.Read(buf)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Printf("%X\n", buf[:numRead])

		ip4header, err := ipv4.ParseHeader(buf[:20])
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		fmt.Println("ipheader:", ip4header)
	}
}

package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal(err)
	}

	path := "/Users/pathbox/tickets.json"
	f, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	info, err := os.Lstat(path)
	if err != nil {
		log.Fatal(err)
		return
	}
	size := info.Size()
	// size := f.Stat().Size()
	n, err := io.CopyN(conn, f, size)
	if err != nil {
		log.Println(err)
		return
	}
	fmt.Println("Pass size: ", n)
}

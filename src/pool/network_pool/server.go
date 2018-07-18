package main

import (
	"bytes"
	"fmt"
	"io"
	"net"
)

func main() {
	addr := "127.0.0.1:9000"

	l, _ := net.Listen("tcp", addr) // tcp socket server
	fmt.Println("Start tcp server: ", addr)
	i := 0
	for {
		conn, _ := l.Accept() // 阻塞接收，外层需要for循环
		br := &bytes.Buffer{}
		n, err := br.ReadFrom(conn)
		i++
		if err != nil {
			if err == io.EOF {
				continue
			} else {
				panic(err)
			}

		}

		fmt.Println("len: ", n)
		fmt.Println("data: ", string(br.Bytes()))
		fmt.Println(i)
	}

}

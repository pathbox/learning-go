package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "127.0.0.1:9091")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Start Listening...")

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		go ParseConn(conn)
	}
}

func ParseConn(conn net.Conn) {
	defer conn.Close()

	// conn.SetDeadline(time.Now().Add(30 * time.Second)) // 30s超时
	var buf bytes.Buffer
	s := time.Now()
	for {
		if _, err := io.CopyN(&buf, conn, 4096); err != nil {
			log.Println("Failed to read record header:", err)
			conn.Close()
			goto OUT
		}
	}

OUT:
	fmt.Println("===Recive: ", len(buf.Bytes()))
	e := time.Now().Sub(s)
	fmt.Println("Need Time: ", e)
	buf.Reset()                   // buf之前所占用的内存不会马上释放，虽然buf清空了，底层还是占用了空间，等待GC释放
	conn.SetDeadline(time.Time{}) // reset deadline
	conn.(*net.TCPConn).SetKeepAlive(true)
	conn.(*net.TCPConn).SetKeepAlivePeriod(3 * time.Minute)
}

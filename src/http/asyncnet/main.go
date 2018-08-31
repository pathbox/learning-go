package main

import (
	"bytes"
	"crypto/tls"
	"io"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	l, err := net.Listen("tcp", "localhost:4242")
	if err != nil {
		log.Fatal(err)
	}

	for {
		conn, err := l.Accept()
		if err, ok := err.(net.Error); ok && err.Temporary() {
			log.Printf("Temporary Accept error: %v; sleeping 1s...", err)
			time.Sleep(1 * time.Second)
			continue
		} else if err != nil {
			log.Fatal(err)
		}
		go serviceConn(conn) // 新起goroutine处理accept得到的conn，避免空的conn连接阻塞main
	}
}

func serviceConn(conn net.Conn) {
	defer conn.Close()

	conn.SetDeadline(time.Now().Add(30 * time.Second)) // 30s超时
	var buf bytes.Buffer
	if _, err := io.CopyN(&buf, conn, 1+2+2); err != nil {
		log.Println("Failed to read record header:", err)
		conn.Close()
		return
	}

	ch, ok := ParseClientHello(buf.Bytes())
	if !ok {
		log.Println("Failed to parse Client Hello.")
	} else {
		log.Printf("Received connection for SNI %q!", ch.SNI)
	}

	conn.SetDeadline(time.Time{}) // reset deadline
	conn.(*net.TCPConn).SetKeepAlive(true)
	conn.(*net.TCPConn).SetKeepAlivePeriod(3 * time.Minute) // 设置KeepAlive

	cert, err := tls.LoadX509KeyPair("localhost.pem", "localhost-key.pem")
	if err != nil {
		log.Fatal(err)
	}
	config := &tls.Config{Certificates: []tls.Certificate{cert}}

	c := tls.Server(prefixConn{
		Reader: io.MultiReader(&buf, conn),
		Conn:   conn,
	}, config)

	proxyConn(c, "gophercon.com:http")
}

type prefixConn struct {
	io.Reader
	net.Conn
}

func (c prefixConn) Read(b []byte) (int, error) {
	return c.Reader.Read(b)
}

func proxyConn(conn net.Conn, addr string) {
	upstream, err := net.Dial("tcp", addr)
	if err != nil {
		log.Println(err)
		return
	}
	defer upstream.Close()
	go io.Copy(upstream, conn) // 将upstream连接请求的数据传给conn
	_, err = io.Copy(conn, upstream)
	log.Printf("Proxy connection finished with err = %v", err)
}

func copyToStderr(conn net.Conn) {
	var total int
	for {
		conn.SetReadDeadline(time.Now().Add(5 * time.Second))
		var buf [128]byte
		n, err := conn.Read(buf[:])
		os.Stderr.Write(buf[:n])
		total += n
		if err != nil {
			log.Printf("Copied %d bytes and ended with err = %v.", total, err)
		}
	}
}

package main

import (
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn) {
	defer c.Close()
	time.Sleep(time.Second * 10)
	for {
		// read from the connection
		time.Sleep(5 * time.Second)
		var buf = make([]byte, 60000) // 建立缓冲buf 类型为 byte的slice，这样每次读取操作(这里是每5秒读一次)，会先从conn连接中读取bytes放到buf中，只读了定义的60000，起到了缓冲的作用
		log.Println("start to read from conn")
		n, err := c.Read(buf)
		if err != nil {
			log.Printf("conn read %d bytes,  error: %s", n, err)
			if nerr, ok := err.(net.Error); ok && nerr.Timeout() {
				continue
			}
			break
		}

		log.Printf("read %d bytes, content is %s\n", n, string(buf[:n]))
	}
}

func main() {
	l, err := net.Listen("tcp", ":9090")
	if err != nil {
		log.Println("listen error: ", err)
		return
	}

	for {
		c, err := l.Accept()
		if err != nil {
			log.Println("accept error: ", err)
			break
		}
		log.Println("accept a new connection")
		go handleConn(c)
	}
}

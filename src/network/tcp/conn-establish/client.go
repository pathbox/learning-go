package main

import (
	"log"
	"net"
	"os"
	"time"
)

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":9090")
	if err != nil {
		log.Printf("%d : dial error: %s", i, err)
		os.Exit(-1)
	}
	log.Println(i, ":connect to server ok")
	return conn
}

func main() {
	var sl []net.Conn
	for i := 1; i < 1000; i++ { //循环创建999 个 conn 放入sl 连接数组中
		conn := establishConn(i)
		if conn != nil {
			sl = append(sl, conn)
		}
	}

	time.Sleep(time.Second * 10000)
}

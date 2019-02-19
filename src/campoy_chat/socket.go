package main

import (
	"io"

	"github.com/golang/net/websocket"
)

type socket struct {
	io.ReadWriterCloser
	done chan bool
}

func (s *socket) Close() error {
	s.done <- true
	return nil
}

func socketHandler(conn *websocket.Conn) { // socket handle 控制入口
	s := socket{conn, make(chan bool)}
	go match(&s) // 每个socket conn 起一个goroutine 处理具体逻辑
	<-s.done     // 可以控制handle是否关闭
}

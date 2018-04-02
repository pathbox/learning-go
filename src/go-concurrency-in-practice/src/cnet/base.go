package ctcp

import (
	"net"
	"time"
)

type TcpMessage struct {
	content string
	err     error
}

// 消息的内容
func (self *TcpMessage) Content() string {
	return self.content
}

func (self *TcpMessage) Err() error {
	return self.err
}

func NewTcpMessage(content string, err error) TcpMessage {
	return TcpMessage{content: content, err: err}
}

// Listener 接口
type TcpListener interface {
	Init(add string) error
	Listen(handler func(conn net.Conn)) error
	Close() bool
	Addr() net.Addr
}

// Sender 接口
type TcpSender interface {
	Init(remoteAddr string, timeout time.Duration) error
	Send(content string) error
	Receive(delim byte) <-chan TcpMessage
	Close() bool
	Addr() net.Addr
	RemoteAddr() net.Addr
}

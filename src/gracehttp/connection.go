package gracehttp

import (
	"net"
)

type Connection struct {
	net.Conn
	listener *Listener
	closed   bool
}

func (conn *Connection) Close() error {
	if !conn.closed {
		conn.closed = true
		conn.listener.wg.Done() //释放wg锁, 等待所有连接都处理完了
	}

	return conn.Conn.Close()
}

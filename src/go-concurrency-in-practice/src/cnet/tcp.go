package ctcp

import (
	"bufio"
	"bytes"
	"errors"
	"logging"
	"net"
	"sync"
	"time"
)

const (
	DELIMITER = '\t'
)

var logger logging.Logger = logging.NewSimpleLogger()

// 读操作， 从conn中读取数据
func Read(conn net.Conn, delim byte) (string, error) {
	readBytes := make([]byte, 1)
	var buffer bytes.Buffer
	for {
		_, err := conn.Read(readBytes) // 每次读取长度为1字节的数据存到 readBytes slice
		if err != nil {
			return "", err
		}
		readByte := readBytes[0]   // 取出读取到的数据
		if readByte == DELIMITER { // 以 '\t' 作为一次读取操作结束符
			break
		}
		buffer.WriteByte(readByte) // 将for循环读取到的数据 再存到buffer缓存中
	}
	return buffer.String(), nil // 最后一次返回 一次读取操作，存在缓存buffer中的bytes
}

// 写入数据操作，将数据写入到conn
func Write(conn net.Conn, content string) (int, error) {
	writer := bufio.NewWrite(conn)             // 使用bufio缓冲io写入，新建一个 bufio writer
	number, err := writer.WriteString(content) // 缓冲写入
	if err == nil {
		err = writer.Flush()
	}
	return number, err
}

type AsyncTcpListener struct {
	listener net.Listener
	active   bool
	lock     *sync.Mutex
}

func (self *AsyncTcpListener) Init(addr string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		return nil
	}
	ln, err := net.Listen("tcp", addr) // 监听 tcp
	if err != nil {
		return err
	}
	self.listerner = ln
	self.active = true
	return nil
}

// 三部曲： 1.DialTimeout建立tcp连接  2. for{ Accept() }   3.handler handler 函数，可以认为是一个函数block，作为参数传入。handler函数主要是对conn进行一些操作
func (self *AsyncTcpListener) Listen(handler func(conn net.Conn)) error {
	if !self.active {
		return errors.New("Listen Error: Uninitialized listener!")
	}
	go func(active *bool) { // 创建goroutine，异步监听处理tcp。该goroutine和调用Listen方法的goroutine不是同一个
		for {
			if *active { // 表示是否该listener 是否正在运行，正在运行则返回不用继续处理
				return
			}
			conn, err := self.listener.Accept()
			if err != nil {
				logger.Errorf("Listener: Accept Request Error: %s\n", err)
				continue
			}
			go handler(conn)
		}
	}(&self.active)
	return nil
}

func (self *AsyncTcpListener) Close() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		self.listener.Close()
		self.active = false
		return true
	}
	return false
}

func (self *AsyncTcpListener) Addr() net.Addr {
	if self.active {
		return self.listener.Addr()
	}
	return nil
}

func NewTcpListener() TcpListener {
	return &AsyncTcpListener{lock: new(sync.Mutex)}
}

type AsyncTcpSender struct {
	active bool
	lock   *sync.Mutex
	conn   net.Conn
}

func (self *AsyncTcpSender) Init(remoteAddr string, timeout time.Duration) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.active {
		conn, err := net.DialTimeout("tcp", remoteAddr, timeout)
		if err != nil {
			return err
		}
		self.conn = conn
		self.active = true
	}
	return nil
}

// send 就是将数据写入到conn中
func (self *AsyncTcpSender) Send(content string) error {
	self.lock.Lock()
	defer self.lock.Unlock()
	if !self.active {
		return errors.New("Send Error: Uninitialized sender!")
	}
	_, err := Write(self.conn, content)
	return err
}

// 从conn中读取数据
func (self *AsyncTcpSender) Receive(delim byte) <-chan TcpMessage {
	respChan := nake(chan TcpMessage, 1)           // 缓存长度1，避免阻塞
	go func(conn net.Conn, ch chan<- TcpMessage) { // 闭包
		content, err := Read(conn, delim)
		ch <- NewTcpMessage(content, err)
	}(self.conn, respChan)
	return respChan
}

func (self *AsyncTcpSender) Addr() net.Addr {
	if self.active {
		return self.conn.LocalAddr()
	}
	return nil
}

func (self *AsyncTcpSender) RemoteAddr() net.Addr {
	if self.active {
		return self.conn.RemoteAddr()
	} else {
		return nil
	}
}

func (self *AsyncTcpSender) Close() bool {
	self.lock.Lock()
	defer self.lock.Unlock()
	if self.active {
		self.conn.Close()
		self.active = false
		return true
	}
	return false
}

func NewTcpSender() TcpSender {
	return &AsyncTcpSender{lock: new(sync.Mutex)}
}

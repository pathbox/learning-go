package main

import (
	"errors"
	"fmt"
	"log"
	"syscall"
)

const TCP_FASTOPEN int = 23

type TFOClient struct {
	ServerAddr [4]byte
	ServerPort int
	fd         int
}

// Create a tcp socket and send data on it. This uses the sendto() system call
// instead of connect() - because connect() calls does not support sending
// data in the syn packet, but the sendto() system call does (as often used in
// connectionless protocols such as udp.
func (c *TFOClient) Send() (err error) {
	c.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0) // 调用Socket接口创建socket,会得到socket地址和socket端口
	if err != nil {
		return
	}
	defer syscall.Close(c.fd)

	sa := &syscall.SockaddrInet4{Addr: c.ServerAddr, Port: c.ServerPort}

	data := []byte("Hello TCP Fast Open")

	log.Printf("Client: Sending to server: %#v\n", string(data))

	// Use the sendto() syscall, instead of connect()
	err = syscall.Sendto(c.fd, data, TCP_FASTOPEN, sa)
	if err != nil {
		if err == syscall.EOPNOTSUPP {
			err = errors.New("TCP Fast Open client support is unavailable (unsupported kernel or disabled, see /proc/sys/net/ipv4/tcp_fastopen).")
		}
		err = errors.New(fmt.Sprintf("Received error in sendTo():", err))
		return
	}

	// Note, this exists before waiting for response and is meant to illustrate
	// the use of the sendto() system call, not of a complete and proper socket
	// setup and teardown processes.

	return
}

package main

import (
	"errors"
	"fmt"
	"log"
	"syscall"
)

type TFOServer struct {
	ServerAddr [4]byte
	ServerPort int
	fd         int
}

const TCP_FASTOPEN int = 23
const LISTEN_BACKLOG int = 23

// Create a tcp socket, setting the TCP_FASTOPEN socket option.
func (s *TFOServer) Bind() (err error) {
	s.fd, err = syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		if err == syscall.ENOPROTOOPT {
			err = errors.New("TCP Fast Open server support is unavailable (unsupported kernel).")
		}
		return
	}
	err = syscall.SetsockoptInt(s.fd, syscall.SOL_TCP, TCP_FASTOPEN, 1)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to set necessary TCP_FASTOPEN socket option: %s", err))
		return
	}

	sa := &syscall.SockaddrInet4{Addr: s.ServerAddr, Port: s.ServerPort}
	err = syscall.Bind(s.fd, sa) // Socket 和地址之间进行绑定
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to bind to Addr: %v, Port: %d, Reason: %s", s.ServerAddr, s.ServerPort, err))
		return
	}

	log.Printf("Server: Bound to addr: %v, port: %d\n", s.ServerAddr, s.ServerPort)

	err = syscall.Listen(s.fd, LISTEN_BACKLOG)
	if err != nil {
		err = errors.New(fmt.Sprintf("Failed to listen: %s", err))
		return
	}

	return
}

// Block, waiting for connections, handling each connection in its own go
// routine
func (s *TFOServer) Accept() {
	log.Println("Server: Waiting for connections")

	defer syscall.Close(s.fd)

	for {
		fd, sockaddr, err := syscall.Accept(s.fd)
		if err != nil {
			log.Fatalln("Failed to accept(): ", err)
		}

		cxn := TFOServerConn{fd: fd, sockaddr: sockaddr.(*syscall.SockaddrInet4)}
		go cxn.Handle()
	}
}

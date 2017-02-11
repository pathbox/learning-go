package socks

import (
	"net"
)

type Socks5Server struct {
	forward Dialer
}

func NewSocks5Server(forward Dialer) (*Socks5Server, error) {
	return &Socks5Server{
		forward: forward,
	}, nil
}

func (s *Socks5Server) Serve(listener net.Listener) error {
	for {
		conn, err := listener.Accept()
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Temporary() {
				continue
			} else {
				return err
			}
		}
		go serveSocks5Client(conn, s.forward)
	}
}

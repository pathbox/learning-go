package web

import (
	"net"
	"net/http/fcgi"
)

func (s *Server) listenAndServeFcgi(addr string) error {
	var l net.Listener
	var err error

	if addr[0] == '/' {
		l, err = net.Listen("unix", addr)
	} else {
		l, err = net.Listen("tcp", addr)
	}

	s.l = l

	if err != nil {
		s.Logger.Println("FCGI listen error", err.Error())
		return err
	}
	return fcgi.Serve(s.l, s)
}

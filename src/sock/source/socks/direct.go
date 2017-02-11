package socks

import (
	"net"
)

type Dialer interface {
	Dial(network, address string) (net.Conn, error)
}

type direct struct{}

var Direct = direct{}

func (direct) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}

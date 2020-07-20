package goseine

import (
	"net"
)

type Dialer interface {
	Dial() (net.Conn, error)
}

type ProxyDialer struct {
	loginServerAddr string
}

func (d *ProxyDialer) Dial() (net.Conn, error) {
	return net.Dial("tcp", d.loginServerAddr)
}

func NewProxyDialer(loginServerAddr string) *ProxyDialer {
	return &ProxyDialer{
		loginServerAddr: loginServerAddr,
	}
}

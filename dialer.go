package goseine

import "net"

type Dialer interface {
	Dial() (*net.TCPConn, error)
}

type ActualAddrResolver interface {
	Resolve() (*net.TCPAddr, error)
}

type ActualRemoteDialer struct {
	resolver ActualAddrResolver
}

func (d *ActualRemoteDialer) Dial() (*net.TCPConn, error) {
	addr, err := d.resolver.Resolve()
	if err != nil {
		return nil, err
	}
	return net.DialTCP("tcp", nil, addr)
}

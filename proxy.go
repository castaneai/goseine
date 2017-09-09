package goseine

import (
	"net"
	"io"
)

const (
	BUFFER_SIZE = 0xFFFF
)

type Proxy struct {
	laddr, raddr *net.TCPAddr
	lconn, rconn io.ReadWriteCloser

	stopCh chan bool
}

func NewProxy(lconn *net.TCPConn, laddr, raddr *net.TCPAddr) (*Proxy, error) {
	return &Proxy {
		lconn: lconn,
		laddr: laddr,
		raddr: raddr,
	}, nil
}

func (p *Proxy) Start() error {
	defer p.lconn.Close()

	var err error
	p.rconn, err = net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		return err
	}
	defer p.rconn.Close()

	go p.pipe(p.lconn, p.rconn)
	go p.pipe(p.rconn, p.lconn)

	// wait for stop
	<- p.stopCh
	return nil
}

func (p *Proxy) pipe(src, dst io.ReadWriter) {
	buff := make([]byte, BUFFER_SIZE)
	for {
		n, err := src.Read(buff)
		if err != nil {
			p.stopCh <- true
			return
		}
		b := buff[:n]

		// TODO: filter

		n, err = dst.Write(b)
		if err != nil {
			p.stopCh <- true
			return
		}
	}
}

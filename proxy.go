package goseine

import (
	"net"
	"io"
	"golang.org/x/sync/errgroup"
	"github.com/sirupsen/logrus"
)

const (
	BUFFER_SIZE = 0xFFFF
)

type ProxyFilter interface {
	Filter(data []byte, src, dst *net.TCPAddr) error
}

type Proxy struct {
	laddr, raddr *net.TCPAddr
	lconn, rconn io.ReadWriteCloser

	log *logrus.Logger
	filter ProxyFilter
}

func NewProxy(lconn *net.TCPConn, laddr, raddr *net.TCPAddr) (*Proxy, error) {
	return &Proxy {
		lconn: lconn,
		laddr: laddr,
		raddr: raddr,
		log: NewLogger("Proxy"),
	}, nil
}

func (p *Proxy) SetFilter(filter ProxyFilter) {
	p.filter = filter
}

func (p *Proxy) Start() error {
	defer (func() {
		p.lconn.Close()
		p.log.Debugf("Local Connection (%s) closed\n", p.laddr.String())
	})()

	var err error
	p.rconn, err = net.DialTCP("tcp", nil, p.raddr)
	if err != nil {
		return err
	}
	defer (func() {
		p.rconn.Close()
		p.log.Debugf("Remote Connection (%s) closed\n", p.raddr.String())
	})()

	eg := &errgroup.Group{}
	eg.Go(func() error {
		lerr := p.pipe(p.lconn, p.rconn, p.laddr, p.raddr)
		if lerr != nil {
			p.log.Debugf("%v\n", lerr)
		}
		return lerr
	})
	eg.Go(func() error {
		rerr := p.pipe(p.rconn, p.lconn, p.raddr, p.laddr)
		if rerr != nil {
			p.log.Debugf("%v\n", rerr)
		}
		return rerr
	})
	return eg.Wait()
}

func (p *Proxy) pipe(src, dst io.ReadWriter, sAddr, dAddr *net.TCPAddr) error {
	buff := make([]byte, BUFFER_SIZE)
	for {
		n, err := src.Read(buff)
		if err != nil {
			return err
		}
		b := buff[:n]

		if p.filter != nil {
			err = p.filter.Filter(b, sAddr, dAddr)
			if err != nil {
				return err
			}
		}

		n, err = dst.Write(b)
		if err != nil {
			return err
		}
	}
}

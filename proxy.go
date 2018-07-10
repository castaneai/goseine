package goseine

import (
	"github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"io"
	"net"
)

type Proxy struct {
	bypass       Dialer
	laddr, raddr *net.TCPAddr
	lconn        io.ReadWriter

	log *logrus.Logger
}

func NewProxy(lconn *net.TCPConn, bypass Dialer) (*Proxy, error) {
	return &Proxy{
		bypass: bypass,
		lconn:  lconn,
		log:    NewLogger("Proxy"),
	}, nil
}

func (p *Proxy) Start() error {
	rconn, err := p.bypass.Dial()
	if err != nil {
		return err
	}
	defer (func() {
		rconn.Close()
		p.log.Debugf("Remote Connection (%s) closed\n", p.raddr.String())
	})()

	eg := &errgroup.Group{}
	eg.Go(func() error {
		lerr := pipe(p.lconn, rconn)
		if lerr != nil {
			p.log.Debugf("%v\n", lerr)
		}
		return lerr
	})
	eg.Go(func() error {
		rerr := pipe(rconn, p.lconn)
		if rerr != nil {
			p.log.Debugf("%v\n", rerr)
		}
		return rerr
	})
	return eg.Wait()
}

func pipe(r io.Reader, w io.Writer) error {
	pr := NewPacketReader(r, NewCurrentGoseineCipher())
	pw := NewPacketWriter(w, NewCurrentGoseineCipher())
	for {
		p := &Packet{}
		if err := pr.Read(p); err != nil {
			return err
		}

		// filtering...

		if err := pw.Write(p); err != nil {
			return err
		}
	}
}

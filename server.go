package goseine

import (
	"crypto/cipher"
	"net"

	"github.com/sirupsen/logrus"
)

type Server struct {
	opts    *serverOptions
	handler PacketHandler
	log     *logrus.Logger
}

type serverOptions struct {
	cipher cipher.Stream
}

var defaultServerOptions = serverOptions{
	cipher: NewDefaultCipher(),
}

type ServerOption interface {
	apply(*serverOptions)
}

type funcServerOption func(*serverOptions)

func (fo funcServerOption) apply(o *serverOptions) {
	fo(o)
}

func NewServer(opt ...ServerOption) *Server {
	opts := defaultServerOptions
	for _, o := range opt {
		o.apply(&opts)
	}
	return &Server{}
}

func (s *Server) Serve(l net.Listener) error {
	s.log.Infof("Listening on %s\n", l.Addr())
	for {
		conn, err := l.Accept()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c net.Conn) {
	defer (func() {
		c.Close()
		s.log.Infof("Closed client: %s\n", c.RemoteAddr().String())
	})()

	s.log.Infof("Handle new client: %s\n", c.RemoteAddr().String())

	w := NewPacketWriter(c, s.opts.cipher)
	r := NewPacketReader(c, s.opts.cipher)
	for {
		var p Packet
		if err := r.Read(&p); err != nil {
			s.log.Error(err)
			return
		}
		s.handler.Handle(w, &p)
	}
}

// http.Handler の真似
type PacketHandler interface {
	// ここでは error を返さない、エラーは log に記録するなり、concrete struct の中でよしなにやってほしい
	Handle(response *PacketWriter, request *Packet)
}

// http.HandlerFunc の真似
type PacketHandlerFunc func(*PacketWriter, *Packet)

func (f PacketHandlerFunc) Handle(response *PacketWriter, request *Packet) {
	f(response, request)
}

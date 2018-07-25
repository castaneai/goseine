package goseine

import (
	"github.com/sirupsen/logrus"
	"net"
)

type Server struct {
	cipher     Cipher
	handler    PacketHandler
	listenAddr *net.TCPAddr
	log        *logrus.Logger
}

func (s *Server) Serve() error {
	lt, err := net.ListenTCP("tcp", s.listenAddr)
	if err != nil {
		return err
	}
	defer lt.Close()

	s.log.Infof("Listening on %s\n", s.listenAddr)
	for {
		conn, err := lt.AcceptTCP()
		if err != nil {
			return err
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(c *net.TCPConn) {
	defer (func() {
		c.Close()
		s.log.Infof("Closed client: %s\n", c.RemoteAddr().String())
	})()

	s.log.Infof("Handle new client: %s\n", c.RemoteAddr().String())

	w := NewPacketWriter(c, s.cipher)
	r := NewPacketReader(c, s.cipher)
	for {
		readPacket := &Packet{}
		if err := r.Read(readPacket); err != nil {
			s.log.Error(err)
			return
		}
		s.handler.Handle(w, readPacket)
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

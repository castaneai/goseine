package goseine

import (
	"github.com/sirupsen/logrus"
	"net"
)

type ActualAddrResolver func() (*net.TCPAddr, error)

type Server struct {
	listenAddr     *net.TCPAddr
	actualResolver ActualAddrResolver
	log            *logrus.Logger
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

func (s *Server) handleConn(c *net.TCPConn) error {
	defer (func() {
		c.Close()
		s.log.Infof("Closed client: %s\n", c.RemoteAddr().String())
	})()

	s.log.Infof("Handle new client: %s\n", c.RemoteAddr().String())

	actual, err := s.actualResolver()
	if err != nil {
		s.log.Errorln(err)
		return err
	}

	p, err := NewProxy(c, c.RemoteAddr().(*net.TCPAddr), actual)
	// f := NewPacketFilter()
	// p.SetFilter(f)

	return p.Start()
}

type LoginServer struct {
	server *Server
}

func NewLoginServer(listenAddr, actualAddr *net.TCPAddr) (*LoginServer, error) {
	log := NewLogger("LoginServer")
	srv := &LoginServer{
		server: &Server{
			listenAddr:     listenAddr,
			actualResolver: func() (*net.TCPAddr, error) { return actualAddr, nil },
			log:            log,
		},
	}
	return srv, nil
}

func Start(listenAddr, actualLoginAddr *net.TCPAddr) error {
	s, err := NewLoginServer(listenAddr, actualLoginAddr)
	if err != nil {
		return err
	}
	return s.server.Serve()
}

package goseine

import (
	"github.com/sirupsen/logrus"
	"net"
)

type Server struct {
	dialer     Dialer
	listenAddr *net.TCPAddr
	log        *logrus.Logger
	proxy      *Proxy
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

	p, err := NewProxy(c, s.dialer)
	if err != nil {
		return err
	}
	s.proxy = p
	// f := NewPacketFilter()
	// p.SetFilter(f)

	return s.proxy.Start()
}

type LoginServer struct {
	server *Server
}

func NewLoginServer(listenAddr *net.TCPAddr, dialer Dialer) (*LoginServer, error) {
	log := NewLogger("LoginServer")
	srv := &LoginServer{
		server: &Server{
			dialer:     dialer,
			listenAddr: listenAddr,
			log:        log,
		},
	}
	return srv, nil
}

func ListenAndServe(listenAddr, actualLoginAddr *net.TCPAddr) error {
	s, err := NewLoginServer(listenAddr, actualLoginAddr)
	if err != nil {
		return err
	}
	return s.server.Serve()
}

package goseine

import (
	"net"
	"testing"
)

func TestAcceptLoginServer(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:10000")
	if err != nil {
		t.Fatal(err)
	}
	actual, err := net.ResolveTCPAddr("tcp", "157.7.172.197:10000")
	if err != nil {
		t.Fatal(err)
	}
	s, err := NewLoginServer(addr, actual)
	if err != nil {
		t.Fatal(err)
	}

	s.server.Serve()
}

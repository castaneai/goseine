package goseine

import (
	"github.com/sirupsen/logrus"
	"testing"
)

func TestRequestPacketLogger(t *testing.T) {
	nopHandler := PacketHandlerFunc(func(w *PacketWriter, req *Packet) {})

	logger := logrus.StandardLogger()
	h := WithMiddleware(nopHandler, RequestPacketLogger(logger))

	p := &Packet{Payload: []byte("\x01\x02\x03\x04"), UseCipher: true}
	h.Handle(nil, p)
}

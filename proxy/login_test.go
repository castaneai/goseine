package proxy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/castaneai/goseine/packets"
)

func TestReplaceChannelServerAddr(t *testing.T) {
	p := &packets.Packet{
		Payload: []byte{
			0x6c, 0xc2, 0xeb, 0x0b, // packet id
			0x00,                                                                                           // status code
			0x31, 0x38, 0x32, 0x2e, 0x32, 0x33, 0x36, 0x2e, 0x36, 0x38, 0x2e, 0x39, 0x30, 0x00, 0x00, 0x00, // ip
			0x23, 0x4e, // port
		},
	}
	l := NewLoginProxy("test", NewPool())
	assert.NoError(t, l.HandleRecv(p))
}

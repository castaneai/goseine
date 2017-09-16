package goseine

import (
	"net"
	"github.com/sirupsen/logrus"
	"encoding/binary"
	"fmt"
)

type PacketFilter struct {
	log *logrus.Logger
	cipher *Cipher
}

func NewPacketFilter() *PacketFilter {
	return &PacketFilter{
		cipher: NewCipher([]byte("qmfaktnpgjs")),
		log: NewLogger("Filter"),
	}
}

func (f *PacketFilter) Filter(data []byte, src, dst *net.TCPAddr) error {
	ptype := "SEND"
	if src.Port >= 10000 && src.Port <= 10099 {
		ptype = "RECV"
	}
	p, err := f.readPacket(data)
	if err != nil {
		return err
	}
	f.log.Infof("[%s]%v\n", ptype, p.Hex())
	return nil
}

func (f *PacketFilter) readPacket(src []byte) (*Packet, error) {
	if len(src) < 4 {
		return nil, fmt.Errorf("Packet source bytes requires least 4 bytes")
	}
	l := binary.LittleEndian.Uint16(src[:2]) - 4
	isCipher := binary.LittleEndian.Uint16(src[2:4]) == 1

	p := make([]byte, l)
	if isCipher {
		f.cipher.Decrypt(p, src[4:4+l])
	} else {
		copy(p, src[4:4+l])
	}
	return NewPacket(p), nil
}
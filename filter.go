package goseine

import (
	"net"
	"github.com/sirupsen/logrus"
	"encoding/binary"
)

type GoseinePacketFilter struct {
	log *logrus.Logger
	cipher *Cipher
}

func (f *GoseinePacketFilter) Filter(data []byte, src, dst *net.TCPAddr) {
	ptype := "SEND"
	if src.Port >= 10000 && src.Port <= 10099 {
		ptype = "RECV"
	}
	p := f.readPacket(data)
	f.log.Infof("[%s]%v\n", ptype, p.Hex())
}

func (f *GoseinePacketFilter) readPacket(src []byte) (*Packet) {
	l := binary.LittleEndian.Uint16(src[:2]) - 4
	isCipher := binary.LittleEndian.Uint16(src[2:4]) == 1

	p := make([]byte, l)
	if isCipher {
		f.cipher.Decrypt(p, src[4:4+l])
	} else {
		copy(p, src[4:4+l])
	}
	return NewPacket(p)
}
package goseine

import (
	"fmt"
	"strings"
)

type Packet struct {
	bytes []byte
	len int
}

func NewPacket(bytes []byte) (*Packet) {
	return &Packet{
		bytes: bytes,
		len: len(bytes),
	}
}

func (p *Packet) Len() int {
	return p.len
}

func (p *Packet) Bytes() []byte {
	return p.bytes
}

func (p *Packet) Hex() string {
	hs := fmt.Sprintf("%X", p.bytes)
	ss := make([]string, p.len, p.len)
	for i := 0; i < p.len; i++ {
		ss[i] = hs[i*2 : i*2+2]
	}
	return strings.Join(ss, " ")
}

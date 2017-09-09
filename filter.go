package goseine

import (
	"net"
	"github.com/sirupsen/logrus"
)

type GoseinePacketFilter struct {
	log *logrus.Logger
}

func (f *GoseinePacketFilter) Filter(data []byte, src, dst *net.TCPAddr) {
	f.log.Infof("[%v -> %v]%v\n", src, dst, data)
}

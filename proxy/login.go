package proxy

import (
	"log"
	"net"

	"github.com/castaneai/goseine/packets"
)

type LoginProxy struct {
	proxy *Proxy
	pool  *Pool
}

func NewLoginProxy(loginServerAddr string, pool *Pool) *LoginProxy {
	proxy := NewProxy(loginServerAddr)
	lp := &LoginProxy{proxy: proxy, pool: pool}
	proxy.FilterReceive(lp)
	return lp
}

func (l *LoginProxy) Start(lis net.Listener) {
	l.proxy.Start(lis)
}

func (l *LoginProxy) HandleRecv(p *packets.Packet) error {
	pid := p.PacketID()
	if pid == packets.PacketIDSelectChannelResponse {
		remoteAddr, err := packets.ReadAddr(p.Payload[5:])
		if err != nil {
			log.Printf("failed to read addr from select channel response: %+v", err)
			return nil
		}
		localAddr, err := l.pool.AssignChannelProxy(remoteAddr)
		if err != nil {
			log.Printf("failed to assign channel proxy: %+v", err)
			return nil
		}
		localAddrBytes, err := packets.WriteAddr(localAddr)
		if err != nil {
			log.Printf("failed to write addr to select channel response: %+v", err)
			return nil
		}
		packets.ReplaceBytes(p.Payload, 5, localAddrBytes)
		log.Printf("replace channel server addr %s -> %s", remoteAddr, localAddr)
	}
	return nil
}

package proxy

import (
	"fmt"
	"net"

	"github.com/castaneai/goseine/packets"
)

type LoginProxy struct {
	proxy *Proxy
	pool  *Pool
}

func NewLoginProxy(loginServerAddr string, pool *Pool) *LoginProxy {
	proxy := NewProxy("LOGINSV", loginServerAddr)
	lp := &LoginProxy{proxy: proxy, pool: pool}
	proxy.FilterReceive(lp)
	return lp
}

func (l *LoginProxy) Start(lis net.Listener) {
	l.proxy.Start(lis)
}

func (l *LoginProxy) HandleRecv(p *packets.Packet) error {
	pid := p.PacketID()
	switch pid {
	case packets.PacketIDSetLoginServerResponse:
		_ = l.replaceSetLoginServerResponse(p)
	case packets.PacketIDSelectChannelResponse:
		_ = l.replaceSelectChannelResponse(p)
	}
	return nil
}

func (l *LoginProxy) replaceSetLoginServerResponse(p *packets.Packet) error {
	beforeAddr, err := packets.ReadAddr(p.Payload[6:])
	if err != nil {
		return fmt.Errorf("failed to read addr from set login server response: %+v", err)
	}
	afterAddr := l.proxy.ListenAddr().String()
	afterAddrBytes, err := packets.WriteAddr(afterAddr)
	if err != nil {
		return fmt.Errorf("failed to write addr: %+v", err)
	}
	packets.ReplaceBytes(p.Payload, 6, afterAddrBytes)
	l.proxy.logger.Debugf("replace set login server %s -> %s", beforeAddr, afterAddr)
	return nil
}

func (l *LoginProxy) replaceSelectChannelResponse(p *packets.Packet) error {
	remoteAddr, err := packets.ReadAddr(p.Payload[5:])
	if err != nil {
		return fmt.Errorf("failed to read addr from select channel response: %+v", err)
	}
	localAddr, err := l.pool.AssignChannelProxy(remoteAddr)
	if err != nil {
		return fmt.Errorf("failed to assign channel proxy: %+v", err)
	}
	localAddrBytes, err := packets.WriteAddr(localAddr)
	if err != nil {
		return fmt.Errorf("failed to write addr to select channel response: %+v", err)
	}
	packets.ReplaceBytes(p.Payload, 5, localAddrBytes)
	l.proxy.logger.Debugf("replace channel server addr %s -> %s", remoteAddr, localAddr)
	return nil
}

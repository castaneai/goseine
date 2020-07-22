package proxy

import (
	"net"
	"sync"
)

type Pool struct {
	mu       sync.RWMutex
	channels map[string]*ChannelProxy
}

func NewPool() *Pool {
	return &Pool{
		mu:       sync.RWMutex{},
		channels: make(map[string]*ChannelProxy),
	}
}

func (p *Pool) AssignChannelProxy(remoteAddr string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if channel, ok := p.channels[remoteAddr]; ok {
		return channel.proxy.lis.Addr().String(), nil
	}

	lis, err := net.Listen("tcp4", "127.0.0.1:0")
	if err != nil {
		return "", err
	}
	proxy := NewProxy(remoteAddr)
	channel := NewChannelProxy(proxy)
	p.channels[remoteAddr] = channel
	go proxy.Start(lis)
	return lis.Addr().String(), nil
}

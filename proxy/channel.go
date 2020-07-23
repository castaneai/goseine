package proxy

import (
	"net"
)

type ChannelProxy struct {
	proxy *Proxy
}

func NewChannelProxy(proxy *Proxy) *ChannelProxy {
	c := &ChannelProxy{proxy: proxy}
	return c
}

func (c *ChannelProxy) Start(lis net.Listener) {
	c.proxy.Start(lis)
}

package proxy

type ChannelProxy struct {
	proxy *Proxy
}

func NewChannelProxy(proxy *Proxy) *ChannelProxy {
	return &ChannelProxy{proxy: proxy}
}

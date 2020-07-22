package proxy

import (
	"net"
	"sync"
)

type Router struct {
	mu    sync.RWMutex
	lrmap map[net.Addr]net.Addr
	rlmap map[net.Addr]net.Addr
}

func (r *Router) Register(laddr, raddr net.Addr) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.lrmap[laddr] = raddr
	r.rlmap[raddr] = laddr
}

func (r *Router) ResolveLocalRemote(addr net.Addr) net.Addr {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if a, ok := r.lrmap[addr]; ok {
		return a
	}
	return addr
}

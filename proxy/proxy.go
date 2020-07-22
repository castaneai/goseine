package proxy

import (
	"encoding/hex"
	"io"
	"log"
	"net"

	"github.com/castaneai/goseine/packets"
)

type SendFilter interface {
	HandleSend(p *packets.Packet) error
}

type RecvFilter interface {
	HandleRecv(p *packets.Packet) error
}

type Proxy struct {
	lis        net.Listener
	remoteAddr string
	sendFilter SendFilter
	recvFilter RecvFilter
}

func NewProxy(remoteAddr string) *Proxy {
	return &Proxy{remoteAddr: remoteAddr}
}

func (p *Proxy) FilterSend(f SendFilter) {
	p.sendFilter = f
}

func (p *Proxy) FilterReceive(f RecvFilter) {
	p.recvFilter = f
}

func (p *Proxy) Start(lis net.Listener) {
	p.lis = lis
	log.Printf("proxy (remote: %v) listening on %v...", p.remoteAddr, lis.Addr())
	for {
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("failed to accept conn: %+v", err)
			return
		}
		// TODO: cancel goroutine
		go p.handleConn(conn)
	}
}

func (p *Proxy) Stop() {
	_ = p.lis.Close()
}

func (p *Proxy) handleConn(lconn net.Conn) {
	log.Printf("new local connected")
	rconn, err := net.Dial("tcp4", p.remoteAddr)
	if err != nil {
		log.Printf("failed to connect to remote: %+v", err)

		return
	}
	log.Printf("new remote connected")
	defer func() { _ = rconn.Close() }()

	// remote -> local
	go func() {
		for {
			var recvPacket packets.Packet
			if err := packets.Read(rconn, &recvPacket); err != nil {
				if err == io.EOF {
					log.Printf("local disconnected (read EOF)")
					return
				}
				log.Printf("failed to read recv packets: %+v", err)
				return
			}

			log.Printf("recv: [%d] %s", recvPacket.PacketID(), hex.Dump(recvPacket.Payload))
			if p.recvFilter != nil {
				if err := p.recvFilter.HandleRecv(&recvPacket); err != nil {
					log.Printf("filter returned error, drop recv packet: %+v", err)
					continue
				}
			}

			if err := packets.Write(lconn, &recvPacket); err != nil {
				log.Printf("failed to write recv packets: %+v", err)
				return
			}
		}
	}()

	// local -> remote
	for {
		var sendPacket packets.Packet
		if err := packets.Read(lconn, &sendPacket); err != nil {
			if err == io.EOF {
				log.Printf("local disconnected (read EOF)")
				return
			}
			log.Printf("failed to read send packets: %+v", err)
			return
		}

		log.Printf("send: [%d] %s", sendPacket.PacketID(), hex.Dump(sendPacket.Payload))
		if p.sendFilter != nil {
			if err := p.sendFilter.HandleSend(&sendPacket); err != nil {
				log.Printf("filter returned error, drop recv packet: %+v", err)
				continue
			}
		}

		if err := packets.Write(rconn, &sendPacket); err != nil {
			log.Printf("failed to write send packets: %+v", err)
			return
		}
	}
}

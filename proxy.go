package goseine

import (
	"encoding/hex"
	"io"
	"log"
	"net"

	"github.com/castaneai/goseine/packet"
)

type Proxy struct {
	dialer Dialer
}

func NewProxy(dialer Dialer) *Proxy {
	return &Proxy{
		dialer: dialer,
	}
}

func (p *Proxy) Start(lis net.Listener) error {
	log.Printf("listening proxy on %v...", lis.Addr())
	for {
		conn, err := lis.Accept()
		if err != nil {
			return err
		}
		go p.handleConn(conn)
	}
}

func (p *Proxy) handleConn(lconn net.Conn) {
	log.Printf("new local connected")
	rconn, err := p.dialer.Dial()
	if err != nil {
		log.Printf("failed to connect to remote: %+v", err)
		return
	}
	log.Printf("new remote connected")
	defer func() { _ = rconn.Close() }()

	// remote -> local
	go func() {
		for {
			var recvPacket packet.Packet
			if err := packet.Read(rconn, &recvPacket); err != nil {
				if err == io.EOF {
					log.Printf("local disconnected (read EOF)")
					return
				}
				log.Printf("failed to read recv packet: %+v", err)
				return
			}

			// TODO: filter
			log.Printf("recv: %s", hex.Dump(recvPacket.Payload))

			if err := packet.Write(lconn, &recvPacket); err != nil {
				log.Printf("failed to write recv packet: %+v", err)
				return
			}
		}
	}()

	// local -> remote
	for {
		var sendPacket packet.Packet
		if err := packet.Read(lconn, &sendPacket); err != nil {
			if err == io.EOF {
				log.Printf("local disconnected (read EOF)")
				return
			}
			log.Printf("failed to read send packet: %+v", err)
			return
		}

		// TODO: filter
		log.Printf("send: %s", hex.Dump(sendPacket.Payload))

		if err := packet.Write(rconn, &sendPacket); err != nil {
			log.Printf("failed to write send packet: %+v", err)
			return
		}
	}
}

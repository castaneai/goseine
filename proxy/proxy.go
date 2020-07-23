package proxy

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/castaneai/goseine/goslog"

	"github.com/sirupsen/logrus"

	"github.com/castaneai/goseine/packets"
)

const (
	packetBufSize = 1024
)

type SendFilter interface {
	HandleSend(p *packets.Packet) error
}

type RecvFilter interface {
	HandleRecv(p *packets.Packet) error
}

type Proxy struct {
	name       string
	logger     logrus.FieldLogger
	lis        net.Listener
	remoteAddr *net.TCPAddr
	sendFilter SendFilter
	recvFilter RecvFilter
}

func NewProxy(name, remoteAddr string) *Proxy {
	ra, err := net.ResolveTCPAddr("tcp4", remoteAddr)
	if err != nil {
		panic(fmt.Sprintf("failed to resolve tcp addr(%s): %+v", remoteAddr, err))
	}
	return &Proxy{
		name:       name,
		logger:     goslog.NewLogger(name),
		remoteAddr: ra,
	}
}

func (p *Proxy) ListenAddr() *net.TCPAddr {
	return p.lis.Addr().(*net.TCPAddr)
}

func (p *Proxy) RemoteAddr() *net.TCPAddr {
	return p.remoteAddr
}

func (p *Proxy) FilterSend(f SendFilter) {
	p.sendFilter = f
}

func (p *Proxy) FilterReceive(f RecvFilter) {
	p.recvFilter = f
}

func (p *Proxy) Start(lis net.Listener) {
	p.lis = lis
	p.logger.Infof("proxy (remote: %s) listening on %s...", p.RemoteAddr(), p.ListenAddr())
	for {
		conn, err := lis.Accept()
		if err != nil {
			p.logger.Errorf("failed to accept conn: %+v", err)
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
	p.logger.Infof("new local connected")
	rconn, err := net.Dial("tcp4", p.RemoteAddr().String())
	if err != nil {
		p.logger.Errorf("failed to connect to remote: %+v", err)
		return
	}
	p.logger.Infof("new remote connected")
	defer func() { _ = rconn.Close() }()

	// remote -> local
	go func() {
		rr := bufio.NewReaderSize(rconn, packetBufSize)
		lw := bufio.NewWriterSize(lconn, packetBufSize)
		for {
			var recvPacket packets.Packet
			if err := packets.Read(rr, &recvPacket); err != nil {
				if err == io.EOF {
					p.logger.Infof("local disconnected (read EOF)")
					return
				}
				p.logger.Errorf("failed to read recv packets: %+v", err)
				return
			}

			if p.recvFilter != nil {
				if err := p.recvFilter.HandleRecv(&recvPacket); err != nil {
					p.logger.Warnf("filter returned error, drop recv packet: %+v", err)
					continue
				}
			}
			p.logger.Debugf(">--- recv: [%d] (len: %d)", recvPacket.PacketID(), len(recvPacket.Payload))

			if err := packets.Write(lw, &recvPacket); err != nil {
				p.logger.Errorf("failed to write recv packets: %+v", err)
				return
			}
			if err := lw.Flush(); err != nil {
				p.logger.Errorf("failed to write recv packets: %+v", err)
				return
			}
		}
	}()

	// local -> remote
	lr := bufio.NewReaderSize(lconn, packetBufSize)
	rw := bufio.NewWriterSize(rconn, packetBufSize)
	for {
		var sendPacket packets.Packet
		if err := packets.Read(lr, &sendPacket); err != nil {
			if err == io.EOF {
				p.logger.Infof("local disconnected (read EOF)")
				return
			}
			p.logger.Errorf("failed to read send packets: %+v", err)
			return
		}

		if sendPacket.PacketID() == packets.PacketIDEnterServerRequest {
			_ = p.replaceEnterServerRequest(&sendPacket)
		}
		if p.sendFilter != nil {
			if err := p.sendFilter.HandleSend(&sendPacket); err != nil {
				p.logger.Warnf("filter returned error, drop recv packet: %+v", err)
				continue
			}
		}
		p.logger.Debugf("<--- send: [%d] (len: %d)", sendPacket.PacketID(), len(sendPacket.Payload))

		if err := packets.Write(rw, &sendPacket); err != nil {
			p.logger.Errorf("failed to write send packets: %+v", err)
			return
		}
		if err := rw.Flush(); err != nil {
			p.logger.Errorf("failed to write send packets: %+v", err)
			return
		}
	}
}

func (p *Proxy) replaceEnterServerRequest(packet *packets.Packet) error {
	beforeIP, err := packets.ReadIPAddr(packet.Payload[12:])
	if err != nil {
		return fmt.Errorf("failed to read IP from enter server request: %+v", err)
	}
	afterIP := p.RemoteAddr().IP.String()
	afterIPBytes, err := packets.WriteIPAddr(afterIP)
	if err != nil {
		return fmt.Errorf("failed to write addr: %+v", err)
	}
	packets.ReplaceBytes(packet.Payload, 12, afterIPBytes)
	p.logger.Debugf("replace enter server request IP: %s -> %s", beforeIP, afterIP)
	return nil
}

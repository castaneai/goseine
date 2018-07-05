package goseine

import (
	"encoding/binary"
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"strings"
)

const (
	packetHeaderLen = 4
)

type Packet struct {
	Payload   []byte
	UseCipher bool
}

type PacketID int

func (p *Packet) PacketID() PacketID {
	return PacketID(binary.LittleEndian.Uint16(p.Payload[:4]))
}

// 01 AB FF C3 のような見やすい16進数表記文字列にする
// ログやデバッグのとき用
func (p *Packet) Hex() string {
	hs := fmt.Sprintf("%X", p.Payload)
	ss := make([]string, len(p.Payload))
	for i := 0; i < len(p.Payload); i++ {
		ss[i] = hs[i*2 : i*2+2]
	}
	return strings.Join(ss, " ")
}

func (p *Packet) String() string {
	return fmt.Sprintf("(ID=%d) %s", p.PacketID(), p.Hex())
}

type packetDecoder struct {
	r      io.Reader
	cipher Cipher
}

func NewPacketDecoder(r io.Reader, cipher Cipher) *packetDecoder {
	return &packetDecoder{r: r, cipher: cipher}
}

func (d *packetDecoder) Decode(dst *Packet) error {
	// header は先頭2byte が長さ, 次の2byteが 暗号化されているかどうかのflag
	header := make([]byte, packetHeaderLen)
	if _, err := io.ReadFull(d.r, header); err != nil {
		return err
	}
	payloadLen := int(binary.LittleEndian.Uint16(header[:2])) - len(header)
	useCipher := binary.LittleEndian.Uint16(header[2:4]) == 1

	dst.UseCipher = useCipher
	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(d.r, payload); err != nil {
		return err
	}
	if useCipher {
		decPayload := make([]byte, payloadLen)
		d.cipher.Decrypt(decPayload, payload)
		dst.Payload = decPayload
	} else {
		dst.Payload = payload
	}
	return nil
}

type packetEncoder struct {
	w      io.Writer
	cipher Cipher
}

func NewPacketEncoder(w io.Writer, cipher Cipher) *packetEncoder {
	return &packetEncoder{w: w, cipher: cipher}
}

func (e *packetEncoder) Encode(src *Packet) error {
	header := make([]byte, packetHeaderLen)
	binary.LittleEndian.PutUint16(header, uint16(len(src.Payload)+packetHeaderLen))
	var useCipherUint16 uint16
	if src.UseCipher {
		useCipherUint16 = 1
	}
	binary.LittleEndian.PutUint16(header[2:], useCipherUint16)
	if _, err := e.w.Write(header); err != nil {
		return err
	}
	if src.UseCipher {
		encPayload := make([]byte, len(src.Payload))
		e.cipher.Encrypt(encPayload, src.Payload)
		if _, err := e.w.Write(encPayload); err != nil {
			return err
		}
	} else {
		if _, err := e.w.Write(src.Payload); err != nil {
			return err
		}
	}
	return nil
}

type PacketWriter interface {
	Write(p *Packet) error
}

// http.Handler の真似
type PacketHandler interface {
	// ここでは error を返さない、エラーは log に記録するなり、concrete struct の中でよしなにやってほしい
	Handle(response PacketWriter, request *Packet)
}

// http.HandlerFunc の真似
type PacketHandlerFunc func(PacketWriter, *Packet)

func (f PacketHandlerFunc) Handle(response PacketWriter, request *Packet) {
	f(response, request)
}

// adapter pattern
// https://medium.com/@matryer/writing-middleware-in-golang-and-how-go-makes-it-so-much-fun-4375c1246e81
type PacketDecorator func(PacketHandler) PacketHandler

func RequestPacketLogger(logger *logrus.Logger) PacketDecorator {
	return func(h PacketHandler) PacketHandler {
		return PacketHandlerFunc(func(w PacketWriter, req *Packet) {
			logger.Println(req)
			h.Handle(w, req)
		})
	}
}

func WithDecorator(h PacketHandler, decorators ...PacketDecorator) PacketHandler {
	for _, decorator := range decorators {
		h = decorator(h)
	}
	return h
}

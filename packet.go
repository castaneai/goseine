package goseine

import (
	"crypto/cipher"
	"encoding/binary"
	"fmt"
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

type PacketReader struct {
	r      io.Reader
	cipher cipher.Stream
}

func NewPacketReader(r io.Reader, cipher cipher.Stream) *PacketReader {
	return &PacketReader{r: r, cipher: cipher}
}

func (pr *PacketReader) Read(dst *Packet) error {
	// header は先頭2byte が長さ, 次の2byteが 暗号化されているかどうかのflag
	header := make([]byte, packetHeaderLen)
	if _, err := io.ReadFull(pr.r, header); err != nil {
		return err
	}
	payloadLen := int(binary.LittleEndian.Uint16(header[:2])) - len(header)
	useCipher := binary.LittleEndian.Uint16(header[2:4]) == 1

	dst.UseCipher = useCipher
	dst.Payload = make([]byte, payloadLen)
	if _, err := io.ReadFull(pr.r, dst.Payload); err != nil {
		return err
	}
	if useCipher {
		pr.cipher.XORKeyStream(dst.Payload, dst.Payload)
	}
	return nil
}

type PacketWriter struct {
	w      io.Writer
	cipher cipher.Stream
}

func NewPacketWriter(w io.Writer, cipher cipher.Stream) *PacketWriter {
	return &PacketWriter{w: w, cipher: cipher}
}

func (pw *PacketWriter) Write(src *Packet) error {
	header := make([]byte, packetHeaderLen)
	binary.LittleEndian.PutUint16(header, uint16(len(src.Payload)+packetHeaderLen))
	var useCipherUint16 uint16
	if src.UseCipher {
		useCipherUint16 = 1
	}
	binary.LittleEndian.PutUint16(header[2:], useCipherUint16)
	if _, err := pw.w.Write(header); err != nil {
		return err
	}
	if src.UseCipher {
		encPayload := make([]byte, len(src.Payload))
		pw.cipher.XORKeyStream(encPayload, src.Payload)
		if _, err := pw.w.Write(encPayload); err != nil {
			return err
		}
	} else {
		if _, err := pw.w.Write(src.Payload); err != nil {
			return err
		}
	}
	return nil
}

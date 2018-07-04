package goseine

import (
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

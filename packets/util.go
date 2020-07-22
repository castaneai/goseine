package packets

import (
	"bytes"
	"fmt"
	"net"
)

func ReadAddr(b []byte) (string, error) {
	if len(b) < 16+2 {
		return "", fmt.Errorf("invalid binary length for addr struct (ip:port)")
	}
	ip := UnpadBytes(b, 16)
	port := byteOrder.Uint16(b[16:])
	return fmt.Sprintf("%s:%d", ip, port), nil
}

func WriteAddr(addr string) ([]byte, error) {
	tcpa, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve \"%s\" to TCP addr: %+v", addr, err)
	}
	ip := PadBytes([]byte(tcpa.IP.String()), 16)
	portb := make([]byte, 2)
	byteOrder.PutUint16(portb, uint16(tcpa.Port))
	return append(ip, portb...), nil
}

func ReplaceBytes(b []byte, pos int, replace []byte) {
	for i := 0; i < len(replace); i++ {
		b[pos+i] = replace[i]
	}
}

func PadBytes(b []byte, size int) []byte {
	if len(b) >= size {
		return b[:size]
	}
	return append(b, bytes.Repeat([]byte{0x00}, size-len(b))...)
}

func UnpadBytes(b []byte, size int) []byte {
	if len(b) >= size {
		return b[:size]
	}
	i := bytes.IndexByte(b, 0x00)
	if i < 0 || i+1 > size {
		return b[:size]
	}
	return b[:i+1]
}

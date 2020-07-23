package packets

import (
	"bytes"
	"fmt"
	"net"
)

func ReadIPAddr(b []byte) (string, error) {
	if len(b) < 16 {
		return "", fmt.Errorf("invalid binary length for addr struct (ip:port)")
	}
	ip := UnpadBytes(b[:16])
	return string(ip), nil
}

func ReadAddr(b []byte) (string, error) {
	ip, err := ReadIPAddr(b[:16])
	if err != nil {
		return "", err
	}
	port := byteOrder.Uint16(b[16:])
	return fmt.Sprintf("%s:%d", ip, port), nil
}

func WriteIPAddr(ip string) ([]byte, error) {
	b := PadBytes([]byte(ip), 16)
	return b, nil
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

func UnpadBytes(b []byte) []byte {
	return bytes.TrimRight(b, "\x00")
}

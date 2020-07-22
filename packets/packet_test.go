package packets

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWrite(t *testing.T) {
	b := new(bytes.Buffer)

	p := &Packet{Payload: []byte("\x01\x02\x03\x04"), UseCipher: true}
	if err := Write(b, p); err != nil {
		t.Fatal(err)
	}

	expected := []byte{0x08, 0x00, 0x01, 0x00, 0x70, 0x6f, 0x65, 0x65}
	assert.Equal(t, expected, b.Bytes())
}

func TestRead(t *testing.T) {
	src := []byte{0x08, 0x00, 0x01, 0x00, 0x70, 0x6f, 0x65, 0x65}

	p := &Packet{UseCipher: true}
	if err := Read(bytes.NewReader(src), p); err != nil {
		t.Fatal(err)
	}

	expected := []byte{0x01, 0x02, 0x03, 0x04}
	assert.Equal(t, expected, p.Payload)
}

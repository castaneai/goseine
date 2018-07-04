package goseine

import (
	"bytes"
	"reflect"
	"testing"
)

func TestPacketEncoder(t *testing.T) {
	actual := new(bytes.Buffer)
	cipher := NewCurrentGoseineCipher()
	enc := NewPacketEncoder(actual, cipher)

	p := &Packet{Payload: []byte("\x01\x02\x03\x04"), UseCipher: true}
	if err := enc.Encode(p); err != nil {
		t.Fatal(err)
	}

	expected := []byte{0x08, 0x00, 0x01, 0x00, 0x70, 0x6f, 0x65, 0x65}
	if !reflect.DeepEqual(expected, actual.Bytes()) {
		t.Fatalf("expected: %#+v, but actual: %#+v", expected, actual.Bytes())
	}
}

func TestPacketDecoder(t *testing.T) {
	srcBytes := []byte{0x08, 0x00, 0x01, 0x00, 0x70, 0x6f, 0x65, 0x65}

	r := bytes.NewReader(srcBytes)
	cipher := NewCurrentGoseineCipher()
	dec := NewPacketDecoder(r, cipher)

	p := &Packet{UseCipher: true}
	if err := dec.Decode(p); err != nil {
		t.Fatal(err)
	}

	expectedPayload := []byte{0x01, 0x02, 0x03, 0x04}
	if !reflect.DeepEqual(expectedPayload, p.Payload) {
		t.Fatalf("expected: %#+v, but actual: %#+v", expectedPayload, p.Payload)
	}
}

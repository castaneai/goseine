package goseine

import (
	"bytes"
	"reflect"
	"testing"
)

func TestPacketWriter(t *testing.T) {
	actual := new(bytes.Buffer)
	cipher := NewDefaultCipher()
	enc := NewPacketWriter(actual, cipher)

	p := &Packet{Payload: []byte("\x01\x02\x03\x04"), UseCipher: true}
	if err := enc.Write(p); err != nil {
		t.Fatal(err)
	}

	expected := []byte{0x08, 0x00, 0x01, 0x00, 0x70, 0x6f, 0x65, 0x65}
	if !reflect.DeepEqual(expected, actual.Bytes()) {
		t.Fatalf("expected: %#+v, but actual: %#+v", expected, actual.Bytes())
	}
}

func TestPacketReader(t *testing.T) {
	srcBytes := []byte{0x08, 0x00, 0x01, 0x00, 0x70, 0x6f, 0x65, 0x65}

	r := bytes.NewReader(srcBytes)
	cipher := NewDefaultCipher()
	dec := NewPacketReader(r, cipher)

	actualPacket := &Packet{UseCipher: true}
	if err := dec.Read(actualPacket); err != nil {
		t.Fatal(err)
	}

	expectedPayload := []byte{0x01, 0x02, 0x03, 0x04}
	if !reflect.DeepEqual(expectedPayload, actualPacket.Payload) {
		t.Fatalf("expected: %#+v, but actual: %#+v", expectedPayload, actualPacket.Payload)
	}
}

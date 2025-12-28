package network

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"testing"

	"frontage/pkg/network/repository"
	"github.com/google/uuid"
)

type testPacket struct {
	Name  string `json:"name"`
	Value int    `json:"value"`
}

func (testPacket) PacketTag() PacketTag { return 0x1234 }

type badPacket struct {
	F func()
}

func (badPacket) PacketTag() PacketTag { return 0x0002 }

func TestSendPacket_NoConnection(t *testing.T) {
	id := uuid.New()
	repository.RemoveConnection(id)

	if SendPacket(id, testPacket{Name: "x", Value: 1}) {
		t.Fatalf("expected false when connection is missing")
	}
}

func TestSendPacket_MarshalError(t *testing.T) {
	id := uuid.New()
	c1, c2 := net.Pipe()
	repository.AddConnection(id, c1)
	t.Cleanup(func() {
		repository.RemoveConnection(id)
		_ = c1.Close()
		_ = c2.Close()
	})

	if SendPacket(id, badPacket{F: func() {}}) {
		t.Fatalf("expected false when json marshal fails")
	}
}

func TestSendPacket_Success(t *testing.T) {
	id := uuid.New()
	c1, c2 := net.Pipe()
	repository.AddConnection(id, c1)
	t.Cleanup(func() {
		repository.RemoveConnection(id)
		_ = c1.Close()
		_ = c2.Close()
	})

	pkt := testPacket{Name: "alpha", Value: 42}
	expectedBody, err := json.Marshal(pkt)
	if err != nil {
		t.Fatalf("marshal test packet: %v", err)
	}

	done := make(chan bool, 1)
	go func() {
		done <- SendPacket(id, pkt)
	}()

	header := make([]byte, 6)
	if _, err := io.ReadFull(c2, header); err != nil {
		t.Fatalf("read header: %v", err)
	}

	tag := binary.LittleEndian.Uint16(header[:2])
	if tag != uint16(pkt.PacketTag()) {
		t.Fatalf("unexpected tag: got %d want %d", tag, pkt.PacketTag())
	}

	bodyLen := binary.LittleEndian.Uint32(header[2:6])
	if bodyLen != uint32(len(expectedBody)) {
		t.Fatalf("unexpected body length: got %d want %d", bodyLen, len(expectedBody))
	}

	body := make([]byte, len(expectedBody))
	if _, err := io.ReadFull(c2, body); err != nil {
		t.Fatalf("read body: %v", err)
	}
	if !bytes.Equal(body, expectedBody) {
		t.Fatalf("unexpected body: got %s want %s", body, expectedBody)
	}

	if ok := <-done; !ok {
		t.Fatalf("expected true on success")
	}
}

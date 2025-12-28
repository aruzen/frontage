package controller

import (
	"bytes"
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"testing"
	"time"

	"frontage/pkg/network"
	"frontage/pkg/network/game_api"
)

// helper to build a packet with header (tag + length) and JSON body
func buildPacket(tag network.PacketTag, body any) ([]byte, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	head := make([]byte, 6)
	binary.LittleEndian.PutUint16(head[:2], uint16(tag))
	binary.LittleEndian.PutUint32(head[2:6], uint32(len(jsonBody)))
	return append(head, jsonBody...), nil
}

// fakeConn is a minimal net.Conn backed by a bytes.Buffer for reads.
type fakeConn struct {
	buf *bytes.Reader
}

func (f *fakeConn) Read(p []byte) (int, error)         { return f.buf.Read(p) }
func (f *fakeConn) Write(p []byte) (int, error)        { return len(p), nil }
func (f *fakeConn) Close() error                       { return nil }
func (f *fakeConn) LocalAddr() net.Addr                { return nil }
func (f *fakeConn) RemoteAddr() net.Addr               { return nil }
func (f *fakeConn) SetDeadline(t time.Time) error      { return nil }
func (f *fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (f *fakeConn) SetWriteDeadline(t time.Time) error { return nil }

// Wrap fakeConn to satisfy *net.TCPConn requirement using interface alias.
type tcpFake struct{ fakeConn }

func TestReceiveLoopRoutesPackets(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	systemCh := make(chan network.Packet, 1)
	lobbyCh := make(chan network.Packet, 1)
	gameCh := make(chan network.Packet, 1)

	// Build one ACT_EVENT packet (Game flag) and one GAME_INITIALIZE packet (Game flag)
	act, err := buildPacket(network.ACT_EVENT_PACKET_TAG, game_api.ActEventPacket{Events: []game_api.ActEventPayload{{}}})
	if err != nil {
		t.Fatal(err)
	}
	init, err := buildPacket(network.GAME_INITIALIZE_PACKET_TAG, game_api.GameInitializePacket{Width: 7, Height: 7, YourSide: 0})
	if err != nil {
		t.Fatal(err)
	}
	data := append(act, init...)

	conn := &tcpFake{fakeConn{buf: bytes.NewReader(data)}}

	errCh := make(chan error, 1)
	go func() {
		errCh <- ReceiveLoop(ctx, conn, systemCh, lobbyCh, gameCh)
	}()

	// consume two packets
	var got []network.Packet
	for i := 0; i < 2; i++ {
		select {
		case p := <-gameCh:
			got = append(got, p)
		case err := <-errCh:
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("ReceiveLoop returned error: %v", err)
			}
		case <-time.After(time.Second):
			t.Fatalf("timeout waiting packet %d", i)
		}
	}

	if len(got) != 2 {
		t.Fatalf("expected 2 packets, got %d", len(got))
	}
	if _, ok := got[0].(game_api.ActEventPacket); !ok {
		t.Fatalf("first packet type mismatch: %T", got[0])
	}
	if _, ok := got[1].(game_api.GameInitializePacket); !ok {
		t.Fatalf("second packet type mismatch: %T", got[1])
	}
}

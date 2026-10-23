package repository

import (
	"encoding/binary"
	"io"
	"net"
	"testing"
	"time"

	"frontage/pkg/network/lobby_api"
	"github.com/google/uuid"
)

func TestSendPacketDeadlineClosesIncompleteFrame(t *testing.T) {
	for _, readHeader := range []bool{false, true} {
		conn, peer := net.Pipe()
		id := uuid.New()
		AddConnection(id, conn)
		t.Cleanup(func() { _ = conn.Close(); _ = peer.Close(); RemoveConnection(id) })
		result := make(chan bool, 1)
		go func() {
			result <- SendPacketWithTimeout(id, lobby_api.CompleteMatchMakePacket{MatchID: uuid.New()}, 100*time.Millisecond)
		}()
		if readHeader {
			var header [6]byte
			_ = peer.SetReadDeadline(time.Now().Add(time.Second))
			if _, err := io.ReadFull(peer, header[:]); err != nil {
				t.Fatal(err)
			}
			if binary.LittleEndian.Uint32(header[2:]) == 0 {
				t.Fatal("empty test body")
			}
		}
		select {
		case sent := <-result:
			if sent {
				t.Fatal("incomplete frame was reported as sent")
			}
		case <-time.After(time.Second):
			t.Fatal("blocked send did not time out")
		}
		_ = peer.SetReadDeadline(time.Now().Add(time.Second))
		if _, err := peer.Read(make([]byte, 1)); err != io.EOF {
			t.Fatalf("failed frame left an open connection: %v", err)
		}
		if SendPacketWithTimeout(id, lobby_api.CompleteMatchMakePacket{}, time.Second) {
			t.Fatal("closed frame stream accepted another packet")
		}
	}
}

func TestSendPacketTimeoutClearsDeadlineAfterSuccess(t *testing.T) {
	conn, peer := net.Pipe()
	_ = peer.SetReadDeadline(time.Now().Add(time.Second))
	id := uuid.New()
	AddConnection(id, conn)
	t.Cleanup(func() { _ = conn.Close(); _ = peer.Close(); RemoveConnection(id) })
	result := make(chan bool, 1)
	go func() {
		result <- SendPacketWithTimeout(id, lobby_api.CompleteMatchMakePacket{MatchID: uuid.New()}, 100*time.Millisecond)
	}()
	var header [6]byte
	if _, err := io.ReadFull(peer, header[:]); err != nil {
		t.Fatal(err)
	}
	body := make([]byte, binary.LittleEndian.Uint32(header[2:]))
	if _, err := io.ReadFull(peer, body); err != nil {
		t.Fatal(err)
	}
	if !<-result {
		t.Fatal("healthy send failed")
	}
	time.Sleep(150 * time.Millisecond)
	_ = peer.SetReadDeadline(time.Now().Add(time.Second))
	go func() { _, err := conn.Write([]byte("x")); result <- err == nil }()
	if _, err := io.ReadFull(peer, make([]byte, 1)); err != nil {
		t.Fatal(err)
	}
	if !<-result {
		t.Fatal("successful send retained its write deadline")
	}
}

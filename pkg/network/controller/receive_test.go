package controller

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"frontage/pkg/network"
	"frontage/pkg/network/repository"
)

func buildRawPacket(tag network.PacketTag, body []byte) []byte {
	header := make([]byte, 6)
	binary.LittleEndian.PutUint16(header[:2], uint16(tag))
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(body)))
	return append(header, body...)
}

func readPacket(t *testing.T, ch <-chan network.UnsolvedPacket, wantTag network.PacketTag, wantBody []byte) {
	t.Helper()
	select {
	case pkt := <-ch:
		if pkt.Tag != wantTag {
			t.Fatalf("unexpected tag: got %d want %d", pkt.Tag, wantTag)
		}
		if string(pkt.Body) != string(wantBody) {
			t.Fatalf("unexpected body: got %q want %q", string(pkt.Body), string(wantBody))
		}
	case <-time.After(2 * time.Second):
		t.Fatalf("timeout waiting for tag %d", wantTag)
	}
}

func TestReceiveLoopRoutesPackets(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = clientConn.Close()
		_ = serverConn.Close()
	}()

	systemCh := make(chan network.UnsolvedPacket, 1)
	lobbyCh := make(chan network.UnsolvedPacket, 1)
	gameBarrier := &repository.BarrierGameChannel{Chan: make(chan network.UnsolvedPacket, 1)}
	gameBarrier.Living.Store(true)

	errCh := make(chan error, 1)
	go func() {
		errCh <- ReceiveLoop(ctx, serverConn, systemCh, lobbyCh, gameBarrier)
	}()

	systemBody := []byte("sys")
	lobbyBody := []byte("lobby")
	gameBody := []byte("game")

	payload := append(buildRawPacket(0, systemBody), buildRawPacket(network.WAIT_MATCH_MAKE_PACKET_TAG, lobbyBody)...)
	payload = append(payload, buildRawPacket(network.GAME_INITIALIZE_PACKET_TAG, gameBody)...)

	if _, err := clientConn.Write(payload); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	readPacket(t, systemCh, 0, systemBody)
	readPacket(t, lobbyCh, network.WAIT_MATCH_MAKE_PACKET_TAG, lobbyBody)
	readPacket(t, gameBarrier.Chan, network.GAME_INITIALIZE_PACKET_TAG, gameBody)

	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
	}
}

func TestReceiveLoopGameBarrier(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverConn, clientConn := net.Pipe()
	defer func() {
		_ = clientConn.Close()
		_ = serverConn.Close()
	}()

	systemCh := make(chan network.UnsolvedPacket, 1)
	lobbyCh := make(chan network.UnsolvedPacket, 1)
	gameBarrier := &repository.BarrierGameChannel{Chan: make(chan network.UnsolvedPacket, 1)}
	gameBarrier.Living.Store(false)

	errCh := make(chan error, 1)
	go func() {
		errCh <- ReceiveLoop(ctx, serverConn, systemCh, lobbyCh, gameBarrier)
	}()

	if _, err := clientConn.Write(buildRawPacket(network.GAME_INITIALIZE_PACKET_TAG, []byte("g1"))); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	select {
	case pkt := <-gameBarrier.Chan:
		t.Fatalf("unexpected game packet: %+v", pkt)
	case <-time.After(300 * time.Millisecond):
	}

	gameBarrier.Living.Store(true)
	if _, err := clientConn.Write(buildRawPacket(network.GAME_INITIALIZE_PACKET_TAG, []byte("g2"))); err != nil {
		t.Fatalf("write payload: %v", err)
	}

	readPacket(t, gameBarrier.Chan, network.GAME_INITIALIZE_PACKET_TAG, []byte("g2"))

	cancel()
	select {
	case <-errCh:
	case <-time.After(time.Second):
	}
}

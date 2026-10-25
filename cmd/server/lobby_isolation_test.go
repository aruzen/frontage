package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"frontage/pkg/network"
	"frontage/pkg/network/data"
	"frontage/pkg/network/lobby_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
)

type blockedLobbyWriteConn struct {
	net.Conn
	started chan struct{}
	once    sync.Once
}

func (conn *blockedLobbyWriteConn) Write(data []byte) (int, error) {
	conn.once.Do(func() { close(conn.started) })
	return conn.Conn.Write(data)
}

func TestLobbySlowPeerDoesNotBlockOtherPeer(t *testing.T) {
	matchRepo = repository.NewMatchRepository()
	cardRepo = repository.NewCardRepository()
	ctx, cancel := context.WithCancel(context.Background())
	visits := make(chan entryAndExitInfo)
	done := make(chan struct{})
	go func() { defer close(done); lobbyLoop(ctx, visits) }()
	slowServer, slowPeer := net.Pipe()
	fastServer, fastPeer := net.Pipe()
	slowID, fastID := uuid.New(), uuid.New()
	slow := &blockedLobbyWriteConn{Conn: slowServer, started: make(chan struct{})}
	repository.AddConnection(slowID, slow)
	repository.AddConnection(fastID, fastServer)
	repository.AddGameChannel(slowID, make(chan network.UnsolvedPacket))
	repository.AddGameChannel(fastID, make(chan network.UnsolvedPacket))
	t.Cleanup(func() {
		cancel()
		_ = slowServer.Close()
		_ = slowPeer.Close()
		_ = fastServer.Close()
		_ = fastPeer.Close()
		repository.RemoveConnection(slowID)
		repository.RemoveConnection(fastID)
		repository.RemoveGameChannel(slowID)
		repository.RemoveGameChannel(fastID)
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Error("lobby did not stop")
		}
	})
	register := func(id uuid.UUID, packets chan network.UnsolvedPacket) {
		t.Helper()
		select {
		case visits <- entryAndExitInfo{id: id, isEntry: true, channel: packets}:
		case <-time.After(time.Second):
			t.Fatal("another peer could not enter the lobby")
		}
	}
	slowPackets := make(chan network.UnsolvedPacket, 1)
	fastPackets := make(chan network.UnsolvedPacket, 1)
	body, err := json.Marshal(lobby_api.WaitMatchMakePacket{Type: data.PvE})
	if err != nil {
		t.Fatal(err)
	}
	request := network.UnsolvedPacket{Tag: network.WAIT_MATCH_MAKE_PACKET_TAG, Body: body}
	register(slowID, slowPackets)
	slowPackets <- request
	select {
	case <-slow.started:
	case <-time.After(time.Second):
		t.Fatal("slow peer did not start writing")
	}
	register(fastID, fastPackets)
	fastPackets <- request
	_ = fastPeer.SetReadDeadline(time.Now().Add(time.Second))
	var header [6]byte
	if _, err := io.ReadFull(fastPeer, header[:]); err != nil {
		t.Fatal(err)
	}
	if tag := binary.LittleEndian.Uint16(header[:2]); tag != uint16(network.COMPLETE_MATCH_MAKE_PACKET_TAG) {
		t.Fatalf("unexpected tag %d", tag)
	}
	size := binary.LittleEndian.Uint32(header[2:])
	if size > 1024 {
		t.Fatalf("unexpected response size %d", size)
	}
	response := make([]byte, size)
	if _, err := io.ReadFull(fastPeer, response); err != nil {
		t.Fatal(err)
	}
	var complete lobby_api.CompleteMatchMakePacket
	if err := json.Unmarshal(response, &complete); err != nil || complete.MatchID == uuid.Nil {
		t.Fatalf("invalid completion packet: %v", err)
	}
}

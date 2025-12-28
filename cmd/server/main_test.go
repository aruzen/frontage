package main

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"io"
	"net"
	"testing"
	"time"

	"frontage/pkg/network"
	"frontage/pkg/network/controller/pve"
	"frontage/pkg/network/data"
	"frontage/pkg/network/game_dispatcher"
	"frontage/pkg/network/game_handler"
	"frontage/pkg/network/lobby_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
)

func TestLobbyMatchMakeSendsCompletePacket(t *testing.T) {
	matchRepo = repository.NewMatchRepository()
	cardRepo = repository.NewCardRepository()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	visitChan := make(chan entryAndExitInfo)
	go lobbyLoop(ctx, visitChan)

	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	id := uuid.New()
	repository.AddConnection(id, serverConn)
	t.Cleanup(func() {
		repository.RemoveConnection(id)
	})

	lobbyChan := make(chan network.UnsolvedPacket, 1)
	visitChan <- entryAndExitInfo{
		id:      id,
		isEntry: true,
		channel: lobbyChan,
	}

	wait := lobby_api.WaitMatchMakePacket{Type: data.PvE}
	payload, err := json.Marshal(wait)
	if err != nil {
		t.Fatalf("marshal wait packet: %v", err)
	}
	lobbyChan <- network.UnsolvedPacket{
		Tag:  network.WAIT_MATCH_MAKE_PACKET_TAG,
		Body: payload,
	}

	_ = clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	header := make([]byte, 6)
	if _, err := io.ReadFull(clientConn, header); err != nil {
		t.Fatalf("read header: %v", err)
	}
	tag := binary.LittleEndian.Uint16(header[:2])
	if network.PacketTag(tag) != network.COMPLETE_MATCH_MAKE_PACKET_TAG {
		t.Fatalf("unexpected tag: got %d want %d", tag, network.COMPLETE_MATCH_MAKE_PACKET_TAG)
	}
	length := binary.LittleEndian.Uint32(header[2:6])
	body := make([]byte, length)
	if _, err := io.ReadFull(clientConn, body); err != nil {
		t.Fatalf("read body: %v", err)
	}
	var resp lobby_api.CompleteMatchMakePacket
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.MatchID == uuid.Nil {
		t.Fatalf("expected non-nil match id")
	}
}

func TestPVEGameSendsActEventPacket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	serverConn, clientConn := net.Pipe()
	t.Cleanup(func() {
		_ = serverConn.Close()
		_ = clientConn.Close()
	})

	id := uuid.New()
	repository.AddConnection(id, serverConn)
	t.Cleanup(func() {
		repository.RemoveConnection(id)
	})
	repository.AddGameChannel(id, make(chan network.UnsolvedPacket, 1))
	t.Cleanup(func() {
		repository.RemoveGameChannel(id)
	})

	rc := pve.RequireContents{
		CardRepo:     repository.NewCardRepository(),
		ActEventDisp: game_dispatcher.NewActEventDispatcher(nil, nil),
		GameInitDisp: game_dispatcher.NewGameInitializeDispatcher(),
	}

	go pve.Game(ctx, rc, id, pve.DefaultGameInfo())

	readPacket := func() (network.PacketTag, []byte) {
		t.Helper()
		_ = clientConn.SetReadDeadline(time.Now().Add(2 * time.Second))
		header := make([]byte, 6)
		if _, err := io.ReadFull(clientConn, header); err != nil {
			t.Fatalf("read header: %v", err)
		}
		tag := network.PacketTag(binary.LittleEndian.Uint16(header[:2]))
		length := binary.LittleEndian.Uint32(header[2:6])
		body := make([]byte, length)
		if _, err := io.ReadFull(clientConn, body); err != nil {
			t.Fatalf("read body: %v", err)
		}
		return tag, body
	}

	tag, _ := readPacket()
	if tag != network.GAME_INITIALIZE_PACKET_TAG {
		t.Fatalf("unexpected init tag: got %d want %d", tag, network.GAME_INITIALIZE_PACKET_TAG)
	}

	tag, body := readPacket()
	if tag != network.ACT_EVENT_PACKET_TAG {
		t.Fatalf("unexpected act event tag: got %d want %d", tag, network.ACT_EVENT_PACKET_TAG)
	}

	handler := game_handler.NewActEventHandler(nil)
	packet, err := handler.ServePacket(body)
	if err != nil {
		t.Fatalf("act event handler error: %v", err)
	}
	if len(packet.Events) == 0 {
		t.Fatalf("expected at least 1 action event")
	}
}

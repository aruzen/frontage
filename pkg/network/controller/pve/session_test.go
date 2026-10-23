package pve

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"frontage/pkg/network"
	"frontage/pkg/network/controller"
	"frontage/pkg/network/game_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
	"io"
	"net"
	"testing"
	"time"
)

func readSyncPacket(t *testing.T, conn net.Conn) game_api.SyncPacket {
	t.Helper()
	conn.SetReadDeadline(time.Now().Add(time.Second))
	var header [6]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		t.Fatal(err)
	}
	if network.PacketTag(binary.LittleEndian.Uint16(header[:2])) != network.SYNC_PACKET_TAG {
		t.Fatal("expected a synchronization request")
	}
	body := make([]byte, binary.LittleEndian.Uint32(header[2:]))
	if _, err := io.ReadFull(conn, body); err != nil {
		t.Fatal(err)
	}
	var packet game_api.SyncPacket
	if err := json.Unmarshal(body, &packet); err != nil {
		t.Fatal(err)
	}
	return packet
}

func TestSyncRequiresMatchingTagAndCause(t *testing.T) {
	for _, scenario := range []string{"retry then acknowledge", "exhausted", "cancelled"} {
		t.Run(scenario, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			server, client := net.Pipe()
			id := uuid.New()
			repository.AddConnection(id, server)
			t.Cleanup(func() { cancel(); server.Close(); client.Close(); repository.RemoveConnection(id) })
			input := make(chan network.UnsolvedPacket, 2)
			game := &gameSession{
				ctx: ctx, id: id, input: input,
				info:    GameInfo{SyncResendMax: 2, SyncResendDuration: 20 * time.Millisecond},
				parsers: controller.MakeDefaultParsers(nil, nil),
			}
			if scenario == "cancelled" {
				game.info.SyncResendDuration = time.Second
			}
			done := make(chan bool, 1)
			go func() { done <- game.syncClient(game_api.SYNC_CAUSE_GAME_START) }()
			request := readSyncPacket(t, client)
			if request.SyncCause != game_api.SYNC_CAUSE_GAME_START {
				t.Fatal("unexpected synchronization cause")
			}
			want := false
			if scenario == "cancelled" {
				cancel()
			} else {
				input <- network.UnsolvedPacket{Tag: network.TURN_PASS_PACKET_TAG, Body: []byte(`{"sync_cause":0}`)}
				input <- network.UnsolvedPacket{Tag: network.SYNC_PACKET_TAG, Body: []byte(`{"sync_cause":1}`)}
				if retry := readSyncPacket(t, client); retry.SyncCause != request.SyncCause {
					t.Fatal("retry changed the synchronization cause")
				}
				if scenario == "retry then acknowledge" {
					input <- network.UnsolvedPacket{Tag: network.SYNC_PACKET_TAG, Body: []byte(`{"sync_cause":0}`)}
					want = true
				}
			}
			select {
			case ok := <-done:
				if ok != want {
					t.Fatalf("synchronization returned %v, want %v", ok, want)
				}
			case <-time.After(time.Second):
				t.Fatal("synchronization did not stop after acknowledgment, exhaustion or cancellation")
			}
		})
	}
}

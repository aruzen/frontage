package main

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"encoding/binary"
	"encoding/json"
	"encoding/pem"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/network"
	"frontage/pkg/network/data"
	"frontage/pkg/network/game_api"
	"frontage/pkg/network/lobby_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
	"io"
	"math/big"
	"net"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func testTLSConfig(t *testing.T) (ServerInfo, *tls.Config) {
	t.Helper()
	public, private, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		t.Fatal(err)
	}
	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		NotBefore:    time.Now().Add(-time.Hour), NotAfter: time.Now().Add(time.Hour),
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageCertSign,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true, IsCA: true,
	}
	certificate, err := x509.CreateCertificate(rand.Reader, template, template, public, private)
	if err != nil {
		t.Fatal(err)
	}
	key, err := x509.MarshalPKCS8PrivateKey(private)
	if err != nil {
		t.Fatal(err)
	}
	dir := t.TempDir()
	info := ServerInfo{
		Addr: "127.0.0.1:0", TLSEnabled: true,
		TLSCertFile: filepath.Join(dir, "cert.pem"), TLSKeyFile: filepath.Join(dir, "key.pem"),
	}
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: certificate})
	if err := os.WriteFile(info.TLSCertFile, certPEM, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(info.TLSKeyFile, pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: key}), 0600); err != nil {
		t.Fatal(err)
	}
	roots := x509.NewCertPool()
	roots.AppendCertsFromPEM(certPEM)
	return info, &tls.Config{RootCAs: roots, MinVersion: tls.VersionTLS12}
}

func TestListenerRejectsMissingTLSCredentials(t *testing.T) {
	for _, info := range []ServerInfo{
		{TLSEnabled: true},
		{TLSEnabled: true, TLSCertFile: "missing.pem"},
		{TLSEnabled: true, TLSCertFile: "missing.pem", TLSKeyFile: "missing-key.pem"},
	} {
		listener, err := newListener(info)
		if err == nil || listener != nil {
			if listener != nil {
				listener.Close()
			}
			t.Fatal("invalid TLS credentials should return an error without a listener")
		}
	}
}

type protocolClient struct {
	t    *testing.T
	conn net.Conn
}

func (c protocolClient) send(packet network.Packet) {
	c.t.Helper()
	body, err := json.Marshal(packet)
	if err != nil {
		c.t.Fatal(err)
	}
	if err := c.conn.SetWriteDeadline(time.Now().Add(3 * time.Second)); err != nil {
		c.t.Fatal(err)
	}
	payload := buildRawPacket(packet.PacketTag(), body)
	for len(payload) > 0 {
		n, err := c.conn.Write(payload)
		if err != nil {
			c.t.Fatal(err)
		}
		payload = payload[n:]
	}
}

func (c protocolClient) read() network.UnsolvedPacket {
	c.t.Helper()
	if err := c.conn.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		c.t.Fatal(err)
	}
	var header [6]byte
	if _, err := io.ReadFull(c.conn, header[:]); err != nil {
		c.t.Fatal(err)
	}
	body := make([]byte, binary.LittleEndian.Uint32(header[2:]))
	if _, err := io.ReadFull(c.conn, body); err != nil {
		c.t.Fatal(err)
	}
	return network.UnsolvedPacket{Tag: network.PacketTag(binary.LittleEndian.Uint16(header[:2])), Body: body}
}

func decodePacket(t *testing.T, body []byte, target interface{}) {
	t.Helper()
	if err := json.Unmarshal(body, target); err != nil {
		t.Fatal(err)
	}
}

func playTestTurns(t *testing.T, client protocolClient, turns int) uuid.UUID {
	t.Helper()
	// A malformed request must not stop the shared lobby loop.
	client.send(lobby_api.WaitMatchMakePacket{Type: data.MatchType(999)})
	client.send(lobby_api.WaitMatchMakePacket{Type: data.PvE})
	first := client.read()
	if first.Tag != network.COMPLETE_MATCH_MAKE_PACKET_TAG {
		t.Fatalf("match confirmation must precede game packets, got %d", first.Tag)
	}
	var match lobby_api.CompleteMatchMakePacket
	decodePacket(t, first.Body, &match)
	if match.MatchID == uuid.Nil {
		t.Fatal("missing match id")
	}
	var playerID uuid.UUID
	initialized, gameStarts, draws := false, 0, 0
	seenCards := make(map[uuid.UUID]bool)
	turn, started := 0, false
	for turn < turns {
		packet := client.read()
		switch packet.Tag {
		case network.SYNC_PACKET_TAG:
			var sync game_api.SyncPacket
			decodePacket(t, packet.Body, &sync)
			client.send(sync)
			if sync.SyncCause == game_api.SYNC_CAUSE_DRAW_PHASE_END && turn%2 == 0 {
				client.send(game_api.TurnPassPacket{Turn: turn})
			}
		case network.GAME_INITIALIZE_PACKET_TAG:
			var init game_api.GameInitializePacket
			decodePacket(t, packet.Body, &init)
			if initialized || init.Width != 7 || init.Height != 7 || init.YourSide != 0 {
				t.Fatalf("unexpected initialization: %+v", init)
			}
			initialized = true
		case network.ACT_EVENT_PACKET_TAG:
			var event game_api.ActEventPacket
			decodePacket(t, packet.Body, &event)
			if len(event.Events) != 1 {
				t.Fatalf("action packet repeated previous effects: got %d effects", len(event.Events))
			}
			switch event.Events[0].ActionTag {
			case string(action.GAME_START_ACTION):
				gameStarts++
			case string(action.PLAYER_DRAW_MAIN_DECK_ACTION), string(action.PLAYER_DRAW_SUB_DECK_ACTION):
				var draw struct {
					PlayerID uuid.UUID `json:"player_id"`
					CardID   uuid.UUID `json:"card_id"`
				}
				encoded, err := json.Marshal(event.Summaries[0][1].Data)
				if err != nil {
					t.Fatal(err)
				}
				decodePacket(t, encoded, &draw)
				if draw.CardID == uuid.Nil || seenCards[draw.CardID] {
					t.Fatal("draw returned a missing or previously drawn card")
				}
				seenCards[draw.CardID] = true
				draws++
				if turn%2 == 0 {
					if playerID != uuid.Nil && playerID != draw.PlayerID {
						t.Fatal("draw used the wrong active player")
					}
					playerID = draw.PlayerID
				} else if playerID == draw.PlayerID {
					t.Fatal("NPC drew from the human player's deck")
				}
			}
		case network.TURN_START_PACKET_TAG:
			var start game_api.TurnStartPacket
			decodePacket(t, packet.Body, &start)
			if !initialized || started || start.Turn != turn || start.IsYourSide != (turn%2 == 0) {
				t.Fatalf("unexpected turn start: %+v, want turn %d", start, turn)
			}
			started = true
			if start.IsYourSide {
				drawType := game_api.DRAW_TYPE_MAIN_ONLY
				if turn == 2 {
					drawType = game_api.DRAW_TYPE_SUB_ONLY
				}
				client.send(game_api.DrawPhasePacket{DrawType: drawType})
			}
		case network.TURN_END_PACKET_TAG:
			var end game_api.TurnEndPacket
			decodePacket(t, packet.Body, &end)
			if !started || end.Turn != turn || end.IsYourSide != (turn%2 == 0) {
				t.Fatalf("unexpected turn end: %+v, want turn %d", end, turn)
			}
			turn++
			started = false
		default:
			t.Fatalf("unexpected packet tag: %d", packet.Tag)
		}
	}
	if gameStarts != 1 || draws != turns || playerID == uuid.Nil {
		t.Fatalf("unexpected game progression: game starts=%d, draws=%d, player=%s", gameStarts, draws, playerID)
	}
	return playerID
}

func waitForGameExit(t *testing.T, channel *repository.BarrierGameChannel, playerID uuid.UUID) {
	t.Helper()
	deadline := time.NewTimer(3 * time.Second)
	defer deadline.Stop()
	tick := time.NewTicker(time.Millisecond)
	defer tick.Stop()
	for {
		_, matched := matchRepo.FindByPlayerID(playerID)
		if !channel.Living.Load() && !matched {
			return
		}
		select {
		case <-tick.C:
		case <-deadline.C:
			t.Fatal("disconnect did not stop the game and remove its match")
		}
	}
}

func TestServerConnectionGameTurnsAndDisconnect(t *testing.T) {
	for _, encrypted := range []bool{false, true} {
		name := "TCP"
		if encrypted {
			name = "TLS"
		}
		t.Run(name, func(t *testing.T) {
			matchRepo = repository.NewMatchRepository()
			cardRepo = repository.NewCardRepository()
			info := ServerInfo{Addr: "127.0.0.1:0"}
			var tlsConfig *tls.Config
			if encrypted {
				info, tlsConfig = testTLSConfig(t)
			}
			listener, err := newListener(info)
			if err != nil {
				t.Fatal(err)
			}
			ctx, cancel := context.WithCancel(context.Background())
			done := make(chan error, 1)
			go func() { done <- serve(ctx, listener) }()
			t.Cleanup(func() {
				cancel()
				select {
				case err := <-done:
					if err != nil {
						t.Error(err)
					}
				case <-time.After(3 * time.Second):
					t.Error("server did not stop after cancellation")
				}
			})
			for reconnect := 0; reconnect < 2; reconnect++ {
				dialer := &net.Dialer{Timeout: 3 * time.Second}
				var conn net.Conn
				if encrypted {
					conn, err = tls.DialWithDialer(dialer, "tcp", listener.Addr().String(), tlsConfig)
				} else {
					conn, err = dialer.Dial("tcp", listener.Addr().String())
				}
				if err != nil {
					t.Fatal(err)
				}
				t.Cleanup(func() { conn.Close() })
				playerID := playTestTurns(t, protocolClient{t: t, conn: conn}, 6)
				channel := repository.GetGameChannel(playerID)
				if channel == nil || !channel.Living.Load() {
					t.Fatal("game was not active")
				}
				conn.Close()
				waitForGameExit(t, channel, playerID)
			}
		})
	}
}

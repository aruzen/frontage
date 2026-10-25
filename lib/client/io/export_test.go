package main

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"frontage/pkg/network"
	"frontage/pkg/network/lobby_api"
)

func installTestClient(t *testing.T, conn net.Conn) {
	t.Helper()
	clientMu.Lock()
	clientConn = conn
	clientMu.Unlock()
	t.Cleanup(func() { CloseClient() })
}

func readMatchMakeFrames(conn net.Conn, count int) error {
	for i := 0; i < count; i++ {
		var header [6]byte
		if _, err := io.ReadFull(conn, header[:]); err != nil {
			return err
		}
		if tag := binary.LittleEndian.Uint16(header[:2]); tag != uint16(network.WAIT_MATCH_MAKE_PACKET_TAG) {
			return fmt.Errorf("frame %d: unexpected tag %d", i, tag)
		}
		size := binary.LittleEndian.Uint32(header[2:])
		if size == 0 || size > 128 {
			return fmt.Errorf("frame %d: invalid body length %d", i, size)
		}
		body := make([]byte, size)
		if _, err := io.ReadFull(conn, body); err != nil {
			return err
		}
		var packet lobby_api.WaitMatchMakePacket
		if err := json.Unmarshal(body, &packet); err != nil {
			return fmt.Errorf("frame %d: %w", i, err)
		}
		if packet.Type != 1 {
			return fmt.Errorf("frame %d: unexpected match type %d", i, packet.Type)
		}
	}
	return nil
}

func TestConcurrentMatchMakeFrames(t *testing.T) {
	conn, peer := net.Pipe()
	installTestClient(t, conn)
	t.Cleanup(func() { _ = peer.Close() })
	if err := peer.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	const senders = 64
	start := make(chan struct{})
	results := make(chan int, senders)
	for i := 0; i < senders; i++ {
		go func() {
			<-start
			results <- int(SendWaitMatchMake(1))
		}()
	}
	close(start)
	err := readMatchMakeFrames(peer, senders)
	_ = peer.Close()
	CloseClient()
	for i := 0; i < senders; i++ {
		if result := <-results; result != 0 && err == nil {
			t.Errorf("send failed: %d", result)
		}
	}
	if err != nil {
		t.Fatal(err)
	}
}

type observedWriteConn struct {
	net.Conn
	writing chan struct{}
	once    sync.Once
}

func (conn *observedWriteConn) Write(data []byte) (int, error) {
	conn.once.Do(func() { close(conn.writing) })
	return conn.Conn.Write(data)
}

func TestCloseClientInterruptsBlockedSend(t *testing.T) {
	conn, peer := net.Pipe()
	observed := &observedWriteConn{Conn: conn, writing: make(chan struct{})}
	installTestClient(t, observed)
	t.Cleanup(func() { _ = peer.Close() })
	sent := make(chan int, 1)
	go func() { sent <- int(SendWaitMatchMake(1)) }()
	select {
	case <-observed.writing:
	case <-time.After(5 * time.Second):
		_ = peer.Close()
		<-sent
		t.Fatal("send did not start")
	}
	closed := make(chan struct{})
	go func() { CloseClient(); close(closed) }()
	closedInTime := false
	select {
	case <-closed:
		closedInTime = true
	case <-time.After(5 * time.Second):
	}
	_ = peer.Close()
	<-closed
	if result := <-sent; result == 0 {
		t.Error("blocked send succeeded after closing its connection")
	}
	if !closedInTime {
		t.Fatal("CloseClient waited for the blocked send")
	}
}

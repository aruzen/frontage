package controller

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/binary"
	"errors"
	"net"
	"testing"
	"time"

	"frontage/pkg/network"
)

func boundedReceiver(t *testing.T, options ReceiveOptions) (net.Conn, <-chan network.UnsolvedPacket, <-chan error, context.CancelFunc) {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	conn, peer := net.Pipe()
	packets := make(chan network.UnsolvedPacket, 1)
	result := make(chan error, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		result <- ReceiveLoopWithOptions(ctx, conn, nil, packets, nil, options)
	}()
	t.Cleanup(func() {
		cancel()
		_ = peer.Close()
		_ = conn.Close()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Error("receiver did not stop")
		}
	})
	return peer, packets, result, cancel
}

func requireReceiveTimeout(t *testing.T, result <-chan error) {
	t.Helper()
	select {
	case err := <-result:
		var timeout net.Error
		if !errors.As(err, &timeout) || !timeout.Timeout() {
			t.Fatalf("expected receive timeout, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("partial frame did not time out")
	}
}

func TestReceiveOptionsBodySizeBoundary(t *testing.T) {
	for _, size := range []int{0, 8} {
		peer, packets, _, _ := boundedReceiver(t, ReceiveOptions{MaxPacketBodySize: 8})
		body := bytes.Repeat([]byte("x"), size)
		if _, err := peer.Write(buildRawPacket(network.WAIT_MATCH_MAKE_PACKET_TAG, body)); err != nil {
			t.Fatal(err)
		}
		readPacket(t, packets, network.WAIT_MATCH_MAKE_PACKET_TAG, body)
	}
	peer, _, result, _ := boundedReceiver(t, ReceiveOptions{MaxPacketBodySize: 8})
	var header [6]byte
	binary.LittleEndian.PutUint16(header[:2], uint16(network.WAIT_MATCH_MAKE_PACKET_TAG))
	binary.LittleEndian.PutUint32(header[2:], 9)
	if _, err := peer.Write(header[:]); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-result:
		if !errors.Is(err, ErrPacketTooLarge) {
			t.Fatalf("expected size error, got %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("oversized body was read")
	}
}

func TestReceiveOptionsPartialFramesExpire(t *testing.T) {
	packet := buildRawPacket(network.WAIT_MATCH_MAKE_PACKET_TAG, []byte("body"))
	for _, prefix := range []int{1, HeaderLength + 1} {
		peer, _, result, _ := boundedReceiver(t, ReceiveOptions{FrameReadTimeout: 100 * time.Millisecond})
		if _, err := peer.Write(packet[:prefix]); err != nil {
			t.Fatal(err)
		}
		requireReceiveTimeout(t, result)
	}
}

func TestReceiveOptionsProgressDoesNotExtendFrameDeadline(t *testing.T) {
	peer, _, result, _ := boundedReceiver(t, ReceiveOptions{FrameReadTimeout: 120 * time.Millisecond})
	packet := buildRawPacket(network.WAIT_MATCH_MAKE_PACKET_TAG, bytes.Repeat([]byte("x"), 8))
	if _, err := peer.Write(packet[:HeaderLength]); err != nil {
		t.Fatal(err)
	}
	sent := make(chan struct{})
	go func() {
		defer close(sent)
		tick := time.NewTicker(30 * time.Millisecond)
		defer tick.Stop()
		for _, value := range packet[HeaderLength:] {
			<-tick.C
			if _, err := peer.Write([]byte{value}); err != nil {
				return
			}
		}
	}()
	t.Cleanup(func() { _ = peer.Close(); <-sent })
	requireReceiveTimeout(t, result)
}

func TestReceiveOptionsAllowsIdleConnectionsAndCancellation(t *testing.T) {
	peer, packets, result, cancel := boundedReceiver(t, ReceiveOptions{FrameReadTimeout: 50 * time.Millisecond})
	select {
	case err := <-result:
		t.Fatalf("idle connection was closed: %v", err)
	case <-time.After(150 * time.Millisecond):
	}
	if _, err := peer.Write(buildRawPacket(network.WAIT_MATCH_MAKE_PACKET_TAG, []byte("ok"))); err != nil {
		t.Fatal(err)
	}
	readPacket(t, packets, network.WAIT_MATCH_MAKE_PACKET_TAG, []byte("ok"))
	cancel()
	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("cancel returned %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("idle read did not accept cancellation")
	}
}

func TestReceiveOptionsBoundsTLSHandshake(t *testing.T) {
	conn, peer := net.Pipe()
	defer peer.Close()
	secure := tls.Server(conn, &tls.Config{GetCertificate: func(*tls.ClientHelloInfo) (*tls.Certificate, error) {
		return nil, errors.New("unexpected client hello")
	}})
	result := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		result <- ReceiveLoopWithOptions(ctx, secure, nil, nil, nil, ReceiveOptions{FrameReadTimeout: 50 * time.Millisecond})
	}()
	select {
	case err := <-result:
		if !errors.Is(err, context.DeadlineExceeded) {
			t.Fatalf("TLS handshake returned %v", err)
		}
	case <-time.After(time.Second):
		_ = conn.Close()
		<-result
		t.Fatal("TLS handshake did not time out")
	}
}

func TestReceiveOptionsRejectsNegativeTimeout(t *testing.T) {
	_, _, result, _ := boundedReceiver(t, ReceiveOptions{FrameReadTimeout: -time.Second})
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("negative timeout was accepted")
		}
	case <-time.After(time.Second):
		t.Fatal("invalid options started reading")
	}
}

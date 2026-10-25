package controller

import (
	"context"
	"encoding/binary"
	"net"
	"testing"
	"time"

	"frontage/pkg/network"
)

func TestReceiveLoopRejectsOversizedHeaderBeforeBody(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	conn, peer := net.Pipe()
	result := make(chan error, 1)
	go func() { result <- ReceiveLoop(ctx, conn, nil, nil, nil) }()
	t.Cleanup(func() { cancel(); _ = peer.Close(); _ = conn.Close() })
	var header [6]byte
	binary.LittleEndian.PutUint16(header[:2], uint16(network.WAIT_MATCH_MAKE_PACKET_TAG))
	binary.LittleEndian.PutUint32(header[2:], (1<<20)+1)
	if _, err := peer.Write(header[:]); err != nil {
		t.Fatal(err)
	}
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("oversized header was accepted")
		}
	case <-time.After(time.Second):
		t.Fatal("receiver waited for an oversized body")
	}
}

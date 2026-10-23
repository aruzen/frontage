package controller

import (
	"context"
	"fmt"
	"net"
	"testing"
	"time"

	"frontage/pkg/network"
	"frontage/pkg/network/repository"
)

func TestReceiveLoopGameStopReleasesPendingPacket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	conn, peer := net.Pipe()
	barrier := &repository.BarrierGameChannel{Chan: make(chan network.UnsolvedPacket)}
	barrier.Living.Store(true)
	lobby := make(chan network.UnsolvedPacket, 1)
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = ReceiveLoop(ctx, conn, nil, lobby, barrier)
	}()
	t.Cleanup(func() {
		cancel()
		_ = peer.Close()
		_ = conn.Close()
		select {
		case <-done:
		case <-barrier.Chan:
			<-done
		}
	})
	if _, err := peer.Write(buildRawPacket(network.GAME_INITIALIZE_PACKET_TAG, []byte("pending"))); err != nil {
		t.Fatal(err)
	}
	second := make(chan error, 1)
	go func() {
		_, err := peer.Write(buildRawPacket(network.WAIT_MATCH_MAKE_PACKET_TAG, []byte("lobby")))
		second <- err
	}()
	select {
	case err := <-second:
		t.Fatalf("game delivery did not block: %v", err)
	case <-time.After(30 * time.Millisecond):
	}
	barrier.Living.Store(false)
	select {
	case packet := <-lobby:
		if string(packet.Body) != "lobby" {
			t.Fatalf("unexpected lobby packet: %q", packet.Body)
		}
	case <-time.After(time.Second):
		t.Fatal("stopping the game did not release packet delivery")
	}
	if err := <-second; err != nil {
		t.Fatal(err)
	}
	_ = peer.Close()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("receive loop survived the disconnect")
	}
}

func TestReceiveLoopCancellationReleasesAllRoutes(t *testing.T) {
	for _, tag := range []network.PacketTag{0, network.WAIT_MATCH_MAKE_PACKET_TAG, network.GAME_INITIALIZE_PACKET_TAG} {
		t.Run(fmt.Sprint(tag), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			conn, peer := net.Pipe()
			system := make(chan network.UnsolvedPacket)
			lobby := make(chan network.UnsolvedPacket)
			barrier := &repository.BarrierGameChannel{Chan: make(chan network.UnsolvedPacket)}
			barrier.Living.Store(true)
			done := make(chan struct{})
			go func() { defer close(done); _ = ReceiveLoop(ctx, conn, system, lobby, barrier) }()
			t.Cleanup(func() {
				cancel()
				_ = peer.Close()
				_ = conn.Close()
				select {
				case <-done:
				case <-system:
					<-done
				case <-lobby:
					<-done
				case <-barrier.Chan:
					<-done
				}
			})
			if _, err := peer.Write(buildRawPacket(tag, []byte("pending"))); err != nil {
				t.Fatal(err)
			}
			cancel()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("cancellation did not release packet delivery")
			}
		})
	}
}

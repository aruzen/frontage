package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"testing"
	"time"

	"frontage/pkg/network"
	"github.com/google/uuid"
)

func TestPacketWorkersKeepOrderAfterInvalidPacket(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	visits := make(chan entryAndExitInfo)
	done := make(chan struct{})
	type received struct {
		id   uuid.UUID
		body string
	}
	results := make(chan received, 4)
	go func() {
		defer close(done)
		packetLoop(ctx, visits, func(_ context.Context, id uuid.UUID, packet network.UnsolvedPacket) error {
			results <- received{id: id, body: string(packet.Body)}
			if string(packet.Body) == "bad" {
				return errors.New("bad packet")
			}
			return nil
		}, nil)
	}()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Error("packet workers did not stop")
		}
	})
	first, second := uuid.New(), uuid.New()
	packets := make(chan network.UnsolvedPacket, 3)
	visits <- entryAndExitInfo{id: first, isEntry: true, channel: packets}
	for _, body := range []string{"bad", "one", "two"} {
		packets <- network.UnsolvedPacket{Body: []byte(body)}
	}
	for _, want := range []string{"bad", "one", "two"} {
		select {
		case got := <-results:
			if got.id != first || got.body != want {
				t.Fatalf("packet order changed: %+v; want %s", got, want)
			}
		case <-time.After(time.Second):
			t.Fatal("invalid packet stopped the connection worker")
		}
	}
	visits <- entryAndExitInfo{id: uuid.New(), isEntry: false}
	visits <- entryAndExitInfo{id: first, isEntry: false}
	otherPackets := make(chan network.UnsolvedPacket, 1)
	visits <- entryAndExitInfo{id: second, isEntry: true, channel: otherPackets}
	otherPackets <- network.UnsolvedPacket{Body: []byte("other")}
	select {
	case got := <-results:
		if got.id != second || got.body != "other" {
			t.Fatalf("wrong connection received packet: %+v", got)
		}
	case <-time.After(time.Second):
		t.Fatal("visitor exit disabled the lobby")
	}
	close(visits)
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("closed visit channel was not handled")
	}
}

func TestReceiveCancellationDuringVisitorNotifications(t *testing.T) {
	for registered := 0; registered <= 2; registered++ {
		t.Run(fmt.Sprint(registered), func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			conn, peer := net.Pipe()
			system, lobby := make(chan entryAndExitInfo), make(chan entryAndExitInfo)
			done := make(chan struct{})
			go func() { defer close(done); receive(ctx, conn, system, lobby) }()
			t.Cleanup(func() {
				cancel()
				_ = peer.Close()
				_ = conn.Close()
				timer := time.NewTimer(time.Second)
				defer timer.Stop()
				for {
					select {
					case <-done:
						return
					case <-system:
					case <-lobby:
					case <-timer.C:
						t.Error("connection cleanup did not finish")
						return
					}
				}
			})
			if registered >= 1 {
				select {
				case <-system:
				case <-time.After(time.Second):
					t.Fatal("system registration did not start")
				}
			}
			if registered == 2 {
				select {
				case <-lobby:
				case <-time.After(time.Second):
					t.Fatal("lobby registration did not start")
				}
			}
			cancel()
			_ = peer.Close()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("receive waited on a stopped visitor loop")
			}
		})
	}
}

func TestPacketWorkerPreservesVisitorContextAndExit(t *testing.T) {
	for _, cancelVisitor := range []bool{true, false} {
		t.Run(fmt.Sprint(cancelVisitor), func(t *testing.T) {
			type contextKey struct{}
			ctx, cancel := context.WithCancel(context.Background())
			visitor, leave := context.WithCancel(context.WithValue(context.Background(), contextKey{}, "visitor"))
			defer cancel()
			defer leave()
			visits := make(chan entryAndExitInfo)
			packets := make(chan network.UnsolvedPacket, 1)
			handled := make(chan context.Context, 1)
			exits := make(chan uuid.UUID, 2)
			done := make(chan struct{})
			go func() {
				defer close(done)
				packetLoop(ctx, visits, func(ctx context.Context, _ uuid.UUID, _ network.UnsolvedPacket) error {
					handled <- ctx
					return nil
				}, func(id uuid.UUID) { exits <- id })
			}()
			id := uuid.New()
			visits <- entryAndExitInfo{id: id, isEntry: true, channel: packets, ctx: visitor}
			packets <- network.UnsolvedPacket{}
			var handlerContext context.Context
			select {
			case handlerContext = <-handled:
				if handlerContext.Value(contextKey{}) != "visitor" {
					t.Fatal("visitor context value was lost")
				}
			case <-time.After(time.Second):
				t.Fatal("packet was not dispatched")
			}
			if cancelVisitor {
				leave()
			} else {
				cancel()
			}
			select {
			case <-handlerContext.Done():
			case <-time.After(time.Second):
				t.Fatal("handler context did not follow its lifetime")
			}
			select {
			case exited := <-exits:
				if exited != id {
					t.Fatalf("wrong visitor exited: %s", exited)
				}
			case <-time.After(time.Second):
				t.Fatal("visitor exit was not reported")
			}
			if cancelVisitor {
				visits <- entryAndExitInfo{id: id}
				close(visits)
			}
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("packet loop did not stop")
			}
			select {
			case <-exits:
				t.Fatal("visitor exit was reported twice")
			default:
			}
		})
	}
}

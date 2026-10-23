package main

import (
	"context"
	"crypto/tls"
	"errors"
	"frontage/pkg/network"
	"frontage/pkg/network/controller"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
	"log/slog"
	"net"
	"reflect"
	"sync"
)

func newListener(info ServerInfo) (net.Listener, error) {
	if info.Addr == "" {
		info.Addr = ":8275"
	}
	if !info.TLSEnabled {
		return net.Listen("tcp", info.Addr)
	}
	if info.TLSCertFile == "" || info.TLSKeyFile == "" {
		return nil, errors.New("TLS requires both a certificate and a private key")
	}
	certificate, err := tls.LoadX509KeyPair(info.TLSCertFile, info.TLSKeyFile)
	if err != nil {
		return nil, err
	}
	return tls.Listen("tcp", info.Addr, &tls.Config{
		Certificates: []tls.Certificate{certificate},
		MinVersion:   tls.VersionTLS12,
	})
}

func serve(ctx context.Context, listener net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)
	stopClose := context.AfterFunc(ctx, func() { _ = listener.Close() })
	var workers sync.WaitGroup
	defer func() {
		cancel()
		_ = listener.Close()
		stopClose()
		workers.Wait()
	}()
	systemVisit := make(chan entryAndExitInfo)
	lobbyVisit := make(chan entryAndExitInfo)
	workers.Add(2)
	go func() { defer workers.Done(); systemLoop(ctx, systemVisit) }()
	go func() { defer workers.Done(); lobbyLoop(ctx, lobbyVisit) }()
	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return err
		}
		workers.Add(1)
		go func() { defer workers.Done(); receive(ctx, conn, systemVisit, lobbyVisit) }()
	}
}

func receive(ctx context.Context, conn net.Conn, systemVisit, lobbyVisit chan entryAndExitInfo) {
	playerCtx, cancel := context.WithCancel(ctx)
	stopClose := context.AfterFunc(ctx, func() { _ = conn.Close() })
	id := uuid.New()
	repository.AddConnection(id, conn)
	systemChan := make(chan network.UnsolvedPacket)
	lobbyChan := make(chan network.UnsolvedPacket)
	gameChannel := repository.AddGameChannel(id, make(chan network.UnsolvedPacket))
	notify := func(ch chan entryAndExitInfo, visit entryAndExitInfo) bool {
		select {
		case ch <- visit:
			return true
		case <-ctx.Done():
			return false
		}
	}
	defer func() {
		_ = conn.Close()
		cancel()
		stopClose()
		notify(systemVisit, entryAndExitInfo{id: id})
		notify(lobbyVisit, entryAndExitInfo{id: id})
		repository.RemoveGameChannel(id)
		repository.RemoveConnection(id)
	}()
	if !notify(systemVisit, entryAndExitInfo{id: id, isEntry: true, channel: systemChan, ctx: playerCtx}) ||
		!notify(lobbyVisit, entryAndExitInfo{id: id, isEntry: true, channel: lobbyChan, ctx: playerCtx}) {
		return
	}
	if err := controller.ReceiveLoop(playerCtx, conn, systemChan, lobbyChan, gameChannel); err != nil &&
		playerCtx.Err() == nil && !errors.Is(err, net.ErrClosed) {
		slog.Debug("Connection closed", "id", id, "err", err)
	}
}

func packetLoop(ctx context.Context, visitChan chan entryAndExitInfo,
	dispatch func(context.Context, uuid.UUID, network.UnsolvedPacket) error, onExit func(uuid.UUID)) {
	cases := []reflect.SelectCase{
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(ctx.Done())},
		{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(visitChan)},
	}
	visitors := make([]entryAndExitInfo, 2)
	remove := func(idx int) {
		if onExit != nil {
			onExit(visitors[idx].id)
		}
		cases = append(cases[:idx], cases[idx+1:]...)
		visitors = append(visitors[:idx], visitors[idx+1:]...)
	}
	defer func() {
		for len(visitors) > 2 {
			remove(len(visitors) - 1)
		}
	}()
	for {
		idx, value, ok := reflect.Select(cases)
		switch idx {
		case 0:
			return
		case 1:
			if !ok {
				return
			}
			visit := value.Interface().(entryAndExitInfo)
			found := -1
			for i := 2; i < len(visitors); i++ {
				if visitors[i].id == visit.id {
					found = i
					break
				}
			}
			if !visit.isEntry {
				if found >= 0 {
					remove(found)
				}
				continue
			}
			if visit.channel == nil || found >= 0 {
				continue
			}
			if visit.ctx == nil {
				visit.ctx = ctx
			}
			cases = append(cases, reflect.SelectCase{Dir: reflect.SelectRecv, Chan: reflect.ValueOf(visit.channel)})
			visitors = append(visitors, visit)
		default:
			if !ok {
				remove(idx)
				continue
			}
			visit := visitors[idx]
			if visit.ctx.Err() != nil {
				continue
			}
			if err := dispatch(visit.ctx, visit.id, value.Interface().(network.UnsolvedPacket)); err != nil {
				slog.Warn("Failed to dispatch packet", "id", visit.id, "err", err)
			}
		}
	}
}

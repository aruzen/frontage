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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	workers := make(map[uuid.UUID]*packetWorker)
	completed := make(chan *packetWorker)
	remove := func(worker *packetWorker) {
		if workers[worker.visit.id] != worker {
			return
		}
		delete(workers, worker.visit.id)
		worker.cancel()
		if onExit != nil {
			onExit(worker.visit.id)
		}
	}
	defer func() {
		for _, worker := range workers {
			remove(worker)
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case worker := <-completed:
			remove(worker)
		case visit, ok := <-visitChan:
			if !ok {
				return
			}
			if !visit.isEntry {
				if worker := workers[visit.id]; worker != nil {
					remove(worker)
				}
				continue
			}
			if visit.channel == nil || workers[visit.id] != nil {
				continue
			}
			if visit.ctx == nil {
				visit.ctx = ctx
			}
			workerCtx, workerCancel := context.WithCancel(visit.ctx)
			worker := &packetWorker{visit: visit, ctx: workerCtx, cancel: workerCancel,
				stop: context.AfterFunc(ctx, workerCancel)}
			workers[visit.id] = worker
			go runPacketWorker(ctx, worker, completed, dispatch)
		}
	}
}

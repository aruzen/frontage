package main

import (
	"context"
	"log/slog"

	"frontage/pkg/network"
	"github.com/google/uuid"
)

type packetWorker struct {
	visit  entryAndExitInfo
	ctx    context.Context
	cancel context.CancelFunc
	stop   func() bool
}

// Each connection processes its existing input channel in order. Cancellation
// follows both the server lifetime and the visitor context passed to the handler.
func runPacketWorker(ctx context.Context, worker *packetWorker, completed chan<- *packetWorker,
	dispatch func(context.Context, uuid.UUID, network.UnsolvedPacket) error) {
	defer func() {
		worker.cancel()
		worker.stop()
		select {
		case completed <- worker:
		case <-ctx.Done():
		}
	}()
	for {
		if worker.ctx.Err() != nil {
			return
		}
		select {
		case <-worker.ctx.Done():
			return
		case packet, ok := <-worker.visit.channel:
			if !ok || worker.ctx.Err() != nil {
				return
			}
			if err := dispatch(worker.ctx, worker.visit.id, packet); err != nil {
				slog.Warn("Failed to dispatch packet", "id", worker.visit.id, "tag", packet.Tag, "err", err)
			}
		}
	}
}

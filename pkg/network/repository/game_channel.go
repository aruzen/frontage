package repository

import (
	"frontage/pkg/network"
	"github.com/google/uuid"
	"log/slog"
	"sync"
	"sync/atomic"
)

var (
	gameChannelTable   = make(map[uuid.UUID]*BarrierGameChannel)
	gameChannelTableMu sync.RWMutex
)

type BarrierGameChannel struct {
	Living atomic.Bool
	Chan   chan network.UnsolvedPacket
}

func GetGameChannel(id uuid.UUID) *BarrierGameChannel {
	gameChannelTableMu.RLock()
	defer gameChannelTableMu.RUnlock()
	gameChan, ok := gameChannelTable[id]
	if !ok {
		slog.Error("GetConnection: connection not found")
	}
	return gameChan
}

func AddGameChannel(id uuid.UUID, gameChan chan network.UnsolvedPacket) *BarrierGameChannel {
	gameChannelTableMu.Lock()
	defer gameChannelTableMu.Unlock()
	v := &BarrierGameChannel{
		Chan: gameChan,
	}
	v.Living.Store(false)
	gameChannelTable[id] = v
	return v
}

func RemoveGameChannel(id uuid.UUID) {
	gameChannelTableMu.Lock()
	defer gameChannelTableMu.Unlock()
	delete(gameChannelTable, id)
}

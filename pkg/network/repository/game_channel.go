package repository

import (
	"frontage/pkg/network"
	"github.com/google/uuid"
	"log/slog"
	"sync"
)

var (
	gameChannelTable   = make(map[uuid.UUID]*BarrierGameChannel)
	gameChannelTableMu sync.RWMutex
)

type BarrierGameChannel struct {
	Living GameChannelStatus
	Chan   chan network.UnsolvedPacket
}

// GameChannelStatus preserves Load/Store while notifying pending deliveries when
// a game stops. Each false-to-true transition starts a new notification lifetime.
// It must not be copied after first use.
type GameChannelStatus struct {
	mu     sync.Mutex
	living bool
	done   chan struct{}
}

func (s *GameChannelStatus) Load() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.living
}

func (s *GameChannelStatus) Store(living bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.living == living {
		return
	}
	if living {
		s.done = make(chan struct{})
	} else if s.done != nil {
		close(s.done)
	}
	s.living = living
}

// Snapshot returns the activity and stop notification for the same game lifetime.
func (s *GameChannelStatus) Snapshot() (bool, <-chan struct{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.done == nil {
		s.done = make(chan struct{})
		if !s.living {
			close(s.done)
		}
	}
	return s.living, s.done
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
	if previous := gameChannelTable[id]; previous != nil {
		previous.Living.Store(false)
	}
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
	if previous := gameChannelTable[id]; previous != nil {
		previous.Living.Store(false)
	}
	delete(gameChannelTable, id)
}

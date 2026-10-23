package repository

import (
	"github.com/google/uuid"
	"sync"
)

type MatchRepository struct {
	mu         sync.RWMutex
	matchIdles []uuid.UUID
	matchIds   []uuid.UUID
	matchTable map[uuid.UUID] /* player */ uuid.UUID /* match */
}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{
		matchIdles: make([]uuid.UUID, 0),
		matchIds:   make([]uuid.UUID, 0),
		matchTable: make(map[uuid.UUID]uuid.UUID),
	}
}

func (m *MatchRepository) Insert(match, p1, p2 uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.matchIds = append(m.matchIds, match)
	m.matchTable[p1] = match
	if p2 != uuid.Nil {
		m.matchTable[p2] = match
	}
}

func (m *MatchRepository) Remove(match uuid.UUID) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, id := range m.matchIds {
		if id == match {
			m.matchIds = append(m.matchIds[:i], m.matchIds[i+1:]...)
		}
	}
	for p, id := range m.matchTable {
		if id == match {
			delete(m.matchTable, p)
		}
	}
}

func (m *MatchRepository) FindByPlayerID(p uuid.UUID) (uuid.UUID, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	match, ok := m.matchTable[p]
	return match, ok
}

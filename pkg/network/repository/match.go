package repository

import "github.com/google/uuid"

type MatchRepository struct {
	matchIdles []uuid.UUID
	matchIds   []uuid.UUID
	matchTable map[uuid.UUID] /* player */ uuid.UUID /* match */
}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{}
}

func (m *MatchRepository) Insert(match, p1, p2 uuid.UUID) {
	m.matchIds = append(m.matchIds, match)
	m.matchTable[p1] = match
	m.matchTable[p2] = match
}

func (m *MatchRepository) Remove(match uuid.UUID) {
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
	match, ok := m.matchTable[p]
	return match, ok
}

package repository

import (
	"frontage/pkg"
	"frontage/pkg/engine/model"
	"sync"
)

type CardRepository struct {
	mu    sync.RWMutex
	table map[pkg.LocalizeTag]model.Card
}

func NewCardRepository() *CardRepository {
	return &CardRepository{
		table: make(map[pkg.LocalizeTag]model.Card),
	}
}

func (r *CardRepository) Insert(card model.Card) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	if _, ok := r.table[card.LocalizeTag()]; ok {
		return ErrDuplicateTag
	}
	r.table[card.LocalizeTag()] = card
	return nil
}

func (r *CardRepository) Find(tag pkg.LocalizeTag) (model.Card, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	card, find := r.table[tag]
	if !find {
		return nil, ErrNotFound
	}
	return card, nil
}

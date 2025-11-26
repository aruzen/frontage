package repository

import (
	"frontage/pkg"
	"frontage/pkg/engine/model"
)

type CardRepository struct {
	table map[pkg.LocalizeTag]model.Card
}

func NewCardRepository() *CardRepository {
	return &CardRepository{
		table: make(map[pkg.LocalizeTag]model.Card),
	}
}

func (r *CardRepository) Insert(card model.Card) error {
	if _, ok := r.table[card.Tag()]; ok {
		return ErrDuplicateTag
	}
	r.table[card.Tag()] = card
	return nil
}

func (r *CardRepository) Find(tag pkg.LocalizeTag) (model.Card, error) {
	card, find := r.table[tag]
	if !find {
		return nil, ErrNotFound
	}
	return card, nil
}

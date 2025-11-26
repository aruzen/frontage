package repository

import "frontage/pkg/engine/logic"

type EffectFinder func(tag logic.EffectActionTag) logic.EffectAction

type EffectRepository struct {
	finder EffectFinder
}

func NewEffectRepository(finder EffectFinder) *EffectRepository {
	return &EffectRepository{finder: finder}
}

func (s *EffectRepository) Find(tag logic.EffectActionTag) (logic.EffectAction, error) {
	f := s.finder(tag)
	if f == nil {
		return nil, ErrNotFound
	}
	return f, nil
}

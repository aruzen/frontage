package repository

import (
	"frontage/pkg"
	"frontage/pkg/engine/logic"
)

type ModifyActionFinder func(tag logic.ModifyActionTag) logic.ModifyAction
type EffectActionFinder func(tag logic.EffectActionTag) logic.EffectAction

type ActionRepository struct {
	modifyFinder ModifyActionFinder
	effectFinder EffectActionFinder
}

func NewEffectRepository(m ModifyActionFinder, e EffectActionFinder) *ActionRepository {
	return &ActionRepository{
		modifyFinder: m,
		effectFinder: e,
	}
}

func (s *ActionRepository) FindEffect(tag logic.EffectActionTag) (logic.EffectAction, error) {
	f := s.effectFinder(tag)
	if f == nil {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *ActionRepository) FindModify(tag logic.ModifyActionTag) (logic.ModifyAction, error) {
	f := s.modifyFinder(tag)
	if f == nil {
		return nil, ErrNotFound
	}
	return f, nil
}

func (s *ActionRepository) Find(tag pkg.LocalizeTag) (logic.Action, error) {
	if a, err := s.FindEffect(logic.EffectActionTag(tag)); err != nil {
		return a, nil
	}
	if a, err := s.FindModify(logic.ModifyActionTag(tag)); err != nil {
		return a, nil
	}
	return nil, ErrNotFound
}

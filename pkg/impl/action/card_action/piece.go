package card_action

import (
	"frontage/internal/event"
	"frontage/pkg/impl/card"
	"frontage/pkg/logic"
	"frontage/pkg/model"
	"log/slog"
)

type PieceActionState struct {
	card  card.Piece
	value int
}

type PieceHPContext struct {
	event.BaseEffectContext
	Value int
}

type PieceMPContext struct {
	event.BaseEffectContext
	Value int
}

type PieceATKContext struct {
	event.BaseEffectContext
	Value int
}

type baseHPAction struct {
	logic.BaseAction[PieceActionState, PieceHPContext]
}

type baseMPAction struct {
	logic.BaseAction[PieceActionState, PieceMPContext]
}

type baseATKAction struct {
	logic.BaseAction[PieceActionState, PieceATKContext]
}

type PieceHPIncreaseAction struct{ baseHPAction }

type PieceHPDecreaseAction struct{ baseHPAction }

type PieceHPFixAction struct{ baseHPAction }

type PieceMPIncreaseAction struct{ baseMPAction }

type PieceMPDecreaseAction struct{ baseMPAction }

type PieceMPFixAction struct{ baseMPAction }

type PieceATKIncreaseAction struct{ baseATKAction }

type PieceATKDecreaseAction struct{ baseATKAction }

type PieceATKFixAction struct{ baseATKAction }

func (p PieceActionState) Card() card.Piece {
	return p.card
}

func (p PieceActionState) Value() int {
	return p.value
}

func (p baseHPAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		return PieceHPContext{event.NewBaseEffectContext(), PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceHPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, hpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetHP(PieceState.card.HP() + hpContext.Value)
	}
}

func (p PieceHPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, hpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetHP(PieceState.card.HP() - hpContext.Value)
	}
}

func (p PieceHPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, hpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetHP(hpContext.Value)
	}
}

func (p baseMPAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		return PieceMPContext{event.NewBaseEffectContext(), PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceMPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, mpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetMP(PieceState.card.MP() + mpContext.Value)
	}
}

func (p PieceMPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, mpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetMP(PieceState.card.MP() - mpContext.Value)
	}
}

func (p PieceMPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, mpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(mpContext.Value)
	}
}

func (p baseATKAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		return PieceATKContext{event.NewBaseEffectContext(), PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceATKIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, atkContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(PieceState.card.ATK() + atkContext.Value)
	}
}

func (p PieceATKDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, atkContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(PieceState.card.ATK() - atkContext.Value)
	}
}

func (p PieceATKFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	PieceState, atkContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(atkContext.Value)
	}
}

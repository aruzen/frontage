package card

import (
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/card"
	"frontage/pkg/event/action"
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
	event.BaseAction[PieceActionState, PieceHPContext]
}

type baseMPAction struct {
	event.BaseAction[PieceActionState, PieceMPContext]
}

type baseATKAction struct {
	event.BaseAction[PieceActionState, PieceATKContext]
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

func (p baseHPAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		return PieceHPContext{event.NewBaseEffectContext(), PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceHPIncreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, hpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetHP(PieceState.card.HP() + hpContext.Value)
	}
}

func (p PieceHPDecreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, hpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetHP(PieceState.card.HP() - hpContext.Value)
	}
}

func (p PieceHPFixAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, hpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetHP(hpContext.Value)
	}
}

func (p baseMPAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		return PieceMPContext{event.NewBaseEffectContext(), PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceMPIncreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, mpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetMP(PieceState.card.MP() + mpContext.Value)
	}
}

func (p PieceMPDecreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, mpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetMP(PieceState.card.MP() - mpContext.Value)
	}
}

func (p PieceMPFixAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, mpContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(mpContext.Value)
	}
}

func (p baseATKAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		return PieceATKContext{event.NewBaseEffectContext(), PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceATKIncreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, atkContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(PieceState.card.ATK() + atkContext.Value)
	}
}

func (p PieceATKDecreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, atkContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(PieceState.card.ATK() - atkContext.Value)
	}
}

func (p PieceATKFixAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	PieceState, atkContext, ok := p.CastStateContext(state, context)
	if ok {
		PieceState.card.SetATK(atkContext.Value)
	}
}

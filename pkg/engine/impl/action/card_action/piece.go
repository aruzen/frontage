package card_action

import (
	"frontage/internal/event"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"log/slog"
)

type PieceActionState struct {
	deckType model.DeckType
	holder   int
	card     card.MutablePiece
	value    int
}

type PieceHPContext struct {
	*event.BaseEffectContext
	Value int
}

type PieceMPContext struct {
	*event.BaseEffectContext
	Value int
}

type PieceATKContext struct {
	*event.BaseEffectContext
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

func (p PieceActionState) HolderIndex() int {
	return p.holder
}

func (p PieceActionState) Card() card.Piece {
	return p.card
}

func (p PieceActionState) Value() int {
	return p.value
}

func updatePieceCard(board *model.Board, state *PieceActionState, mutate func(card.MutablePiece)) *model.Board {
	if board == nil || state == nil || mutate == nil {
		return board
	}

	next := board.Next()
	cloned := state.card.Copy()
	mutate(cloned)

	players := next.Players()
	if state.holder < 0 || state.holder >= len(players) {
		slog.Warn("invalid holder index for piece update.", "holder", state.holder)
		return next
	}

	player := players[state.holder]
	deckPlayer, ok := player.(model.DeckPlayer)
	if !ok {
		slog.Warn("player does not support deck operations.")
		return next
	}
	deck := deckPlayer.GetDeck(state.deckType)
	if deck == nil {
		slog.Warn("deck is nil for holder", "holder", state.holder, "deckType", state.deckType)
		return next
	}

	if !deck.Update(cloned) {
		slog.Warn("failed to update card in deck.", "deckType", state.deckType, "cardID", cloned.Id())
	}

	return next
}

func (p baseHPAction) MakeState(holder int, deckType model.DeckType, card card.MutablePiece, value int) PieceActionState {
	return PieceActionState{
		deckType: deckType,
		holder:   holder,
		card:     card,
		value:    value,
	}
}

func (p baseMPAction) MakeState(holder int, deckType model.DeckType, card card.MutablePiece, value int) PieceActionState {
	return PieceActionState{
		deckType: deckType,
		holder:   holder,
		card:     card,
		value:    value,
	}
}

func (p baseATKAction) MakeState(holder int, deckType model.DeckType, card card.MutablePiece, value int) PieceActionState {
	return PieceActionState{
		deckType: deckType,
		holder:   holder,
		card:     card,
		value:    value,
	}
}

func (p baseHPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		base := event.NewBaseEffectContext()
		return PieceHPContext{BaseEffectContext: &base, Value: PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceHPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, hpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceHPIncreaseAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetHP(pieceState.card.HP() + hpContext.Value)
	})
}

func (p PieceHPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, hpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceHPDecreaseAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetHP(pieceState.card.HP() - hpContext.Value)
	})
}

func (p PieceHPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, hpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceHPFixAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetHP(hpContext.Value)
	})
}

func (p baseMPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		base := event.NewBaseEffectContext()
		return PieceMPContext{BaseEffectContext: &base, Value: PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceMPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, mpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceMPIncreaseAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetMP(pieceState.card.MP() + mpContext.Value)
	})
}

func (p PieceMPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, mpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceMPDecreaseAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetMP(pieceState.card.MP() - mpContext.Value)
	})
}

func (p PieceMPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, mpContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceMPFixAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetMP(mpContext.Value)
	})
}

func (p baseATKAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if PieceState, ok := state.(PieceActionState); ok {
		base := event.NewBaseEffectContext()
		return PieceATKContext{BaseEffectContext: &base, Value: PieceState.value}
	}
	slog.Warn("State was not PieceActionState.")
	return nil
}

func (p PieceATKIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, atkContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceATKIncreaseAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetATK(pieceState.card.ATK() + atkContext.Value)
	})
}

func (p PieceATKDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, atkContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceATKDecreaseAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetATK(pieceState.card.ATK() - atkContext.Value)
	})
}

func (p PieceATKFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	pieceState, atkContext, ok := p.CastStateContext(state, context)
	if !ok {
		slog.Warn("CastStateContext failed in PieceATKFixAction.")
		return board
	}

	return updatePieceCard(board, pieceState, func(m card.MutablePiece) {
		m.SetATK(atkContext.Value)
	})
}

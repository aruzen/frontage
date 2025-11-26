package piece_action

import (
	"frontage/internal/event"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"log/slog"
)

type PieceOperateActionState struct {
	piece model.MutablePiece
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

type basePieceHPAction struct {
	logic.BaseAction[PieceOperateActionState, PieceHPContext]
}

type basePieceMPAction struct {
	logic.BaseAction[PieceOperateActionState, PieceMPContext]
}

type basePieceATKAction struct {
	logic.BaseAction[PieceOperateActionState, PieceATKContext]
}

type PieceHPIncreaseAction struct{ basePieceHPAction }

type PieceHPDecreaseAction struct{ basePieceHPAction }

type PieceHPFixAction struct{ basePieceHPAction }

type PieceMPIncreaseAction struct{ basePieceMPAction }

type PieceMPDecreaseAction struct{ basePieceMPAction }

type PieceMPFixAction struct{ basePieceMPAction }

type PieceATKIncreaseAction struct{ basePieceATKAction }

type PieceATKDecreaseAction struct{ basePieceATKAction }

type PieceATKFixAction struct{ basePieceATKAction }

func (e PieceOperateActionState) Piece() model.MutablePiece {
	return e.piece
}

func (e PieceOperateActionState) Value() int {
	return e.value
}

func (b basePieceHPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if pieceState, ok := state.(PieceOperateActionState); ok {
		return &PieceHPContext{event.NewBaseEffectContext(), pieceState.value}
	}
	slog.Warn("State was not PieceOperateActionState.")
	return nil
}

func (e PieceHPIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetHP(pieceState.piece.HP() + context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (e PieceHPDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetHP(pieceState.piece.HP() - context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (e PieceHPFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetHP(context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (b basePieceMPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if pieceState, ok := state.(PieceOperateActionState); ok {
		return &PieceMPContext{event.NewBaseEffectContext(), pieceState.value}
	}
	slog.Warn("State was not PieceOperateActionState.")
	return nil
}

func (e PieceMPIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetMP(pieceState.piece.MP() + context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (e PieceMPDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetMP(pieceState.piece.MP() - context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (e PieceMPFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetMP(context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (b basePieceATKAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if pieceState, ok := state.(PieceOperateActionState); ok {
		return &PieceATKContext{event.NewBaseEffectContext(), pieceState.value}
	}
	slog.Warn("State was not PieceOperateActionState.")
	return nil
}

func (e PieceATKIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetATK(pieceState.piece.ATK() + context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (e PieceATKDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetATK(pieceState.piece.ATK() - context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

func (e PieceATKFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	pieceState.piece.SetATK(context.Value)
	board.UpdatePiece(pieceState.piece)
	return board
}

package piece_action

import (
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
)

type PieceSummonActionState struct {
	point pkg.Point
	piece model.Piece
}

type PieceMoveActionState struct {
	from, to pkg.Point
	piece    model.Piece
}

type PieceAttackActionState struct {
	decreaseHPState *PieceOperateActionState
	point           pkg.Point
	piece           model.Piece
	value           int
}

type PieceSummonActionContext struct {
	event.BaseEffectContext
	Point pkg.Point
	Piece model.Piece
}

type PieceMoveActionContext struct {
	event.BaseEffectContext
	Point pkg.Point
}

type PieceAttackActionContext struct {
	event.BaseEffectContext
	Point pkg.Point
	Value int
}

type PieceSummonAction struct {
	logic.BaseAction[PieceSummonActionState, PieceSummonActionContext]
}

type PieceMoveAction struct {
	logic.BaseAction[PieceMoveActionState, PieceMoveActionContext]
}

type PieceAttackAction struct {
	logic.BaseAction[PieceAttackActionState, PieceAttackActionContext]
}

type PieceInvasionAction struct {
	logic.BaseAction[PieceAttackActionState, PieceAttackActionContext]
}

func (e PieceSummonAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(PieceSummonActionState); ok {
		return &PieceSummonActionContext{event.BaseEffectContext{}, s.point, s.piece}
	}
	return nil
}

func (e PieceSummonAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	_, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()
	if !board.SetPiece(c.Point, c.Piece) {
		return nil
	}
	return board
}

func (e PieceMoveAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(PieceMoveActionState); ok {
		return &PieceSummonActionContext{event.BaseEffectContext{}, s.to, s.piece}
	}
	return nil
}

func (e PieceMoveAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()

	if !board.SetPiece(c.Point, s.piece) {
		return nil
	}
	return board
}

func (e PieceAttackAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(PieceMoveActionState); ok {
		return &PieceSummonActionContext{event.BaseEffectContext{}, s.to, s.piece}
	}
	return nil
}

func (e PieceAttackAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()
	s.decreaseHPState.piece = board.Entities()[c.Point.X][c.Point.Y].Copy()
	s.decreaseHPState.value = c.Value
	return board
}

func (e PieceAttackAction) SubEffects(state interface{}) []*logic.EffectEvent {
	s, ok := state.(PieceAttackActionState)
	if !ok {
		return nil
	}
	result := make([]*logic.EffectEvent, 1)
	result[0] = logic.NewEffectEvent(action.FindActionEffect(action.ENTITY_HP_DECREASE_ACTION), s.decreaseHPState)
	return result
}

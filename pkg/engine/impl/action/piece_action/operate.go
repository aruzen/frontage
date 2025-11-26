package piece_action

import (
	"fmt"
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
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

func (c PieceHPContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["value"] = c.Value
	return result
}

func (c *PieceHPContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["value"]; ok {
		num, err := pkg.ToInt(v)
		if err != nil {
			return fmt.Errorf("value: %w", err)
		}
		c.Value = num
	}
	return nil
}

func (c PieceMPContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["value"] = c.Value
	return result
}

func (c *PieceMPContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["value"]; ok {
		num, err := pkg.ToInt(v)
		if err != nil {
			return fmt.Errorf("value: %w", err)
		}
		c.Value = num
	}
	return nil
}

func (c PieceATKContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["value"] = c.Value
	return result
}

func (c *PieceATKContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["value"]; ok {
		num, err := pkg.ToInt(v)
		if err != nil {
			return fmt.Errorf("value: %w", err)
		}
		c.Value = num
	}
	return nil
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

func (PieceHPIncreaseAction) Tag() logic.EffectActionTag      { return action.ENTITY_HP_INCREASE_ACTION }
func (PieceHPDecreaseAction) Tag() logic.EffectActionTag      { return action.ENTITY_HP_DECREASE_ACTION }
func (PieceHPFixAction) Tag() logic.EffectActionTag           { return action.ENTITY_HP_FIX_ACTION }
func (PieceMPIncreaseAction) Tag() logic.EffectActionTag      { return action.ENTITY_MP_INCREASE_ACTION }
func (PieceMPDecreaseAction) Tag() logic.EffectActionTag      { return action.ENTITY_MP_DECREASE_ACTION }
func (PieceMPFixAction) Tag() logic.EffectActionTag           { return action.ENTITY_MP_FIX_ACTION }
func (PieceATKIncreaseAction) Tag() logic.EffectActionTag     { return action.ENTITY_ATK_INCREASE_ACTION }
func (PieceATKDecreaseAction) Tag() logic.EffectActionTag     { return action.ENTITY_ATK_DECREASE_ACTION }
func (PieceATKFixAction) Tag() logic.EffectActionTag          { return action.ENTITY_ATK_FIX_ACTION }
func (a PieceHPIncreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceHPDecreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceHPFixAction) LocalizeTag() pkg.LocalizeTag       { return pkg.LocalizeTag(a.Tag()) }
func (a PieceMPIncreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceMPDecreaseAction) LocalizeTag() pkg.LocalizeTag  { return pkg.LocalizeTag(a.Tag()) }
func (a PieceMPFixAction) LocalizeTag() pkg.LocalizeTag       { return pkg.LocalizeTag(a.Tag()) }
func (a PieceATKIncreaseAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }
func (a PieceATKDecreaseAction) LocalizeTag() pkg.LocalizeTag { return pkg.LocalizeTag(a.Tag()) }
func (a PieceATKFixAction) LocalizeTag() pkg.LocalizeTag      { return pkg.LocalizeTag(a.Tag()) }

func (e PieceOperateActionState) Piece() model.MutablePiece {
	return e.piece
}

func (e PieceOperateActionState) Value() int {
	return e.value
}

func (b basePieceHPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if pieceState, ok := state.(PieceOperateActionState); ok {
		return &PieceHPContext{event.NewBaseEffectContext(), pieceState.value}, logic.Summary{"value": pieceState.value}
	}
	slog.Warn("State was not PieceOperateActionState.")
	return nil, nil
}

func (e PieceHPIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetHP(pieceState.piece.HP() + context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"delta_hp": context.Value}
}

func (e PieceHPDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetHP(pieceState.piece.HP() - context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"delta_hp": -context.Value}
}

func (e PieceHPFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetHP(context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"set_hp": context.Value}
}

func (b basePieceMPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if pieceState, ok := state.(PieceOperateActionState); ok {
		return &PieceMPContext{event.NewBaseEffectContext(), pieceState.value}, logic.Summary{"value": pieceState.value}
	}
	slog.Warn("State was not PieceOperateActionState.")
	return nil, nil
}

func (e PieceMPIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetMP(pieceState.piece.MP() + context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"delta_mp": context.Value}
}

func (e PieceMPDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetMP(pieceState.piece.MP() - context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"delta_mp": -context.Value}
}

func (e PieceMPFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetMP(context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"set_mp": context.Value}
}

func (b basePieceATKAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if pieceState, ok := state.(PieceOperateActionState); ok {
		return &PieceATKContext{event.NewBaseEffectContext(), pieceState.value}, logic.Summary{"value": pieceState.value}
	}
	slog.Warn("State was not PieceOperateActionState.")
	return nil, nil
}

func (e PieceATKIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetATK(pieceState.piece.ATK() + context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"delta_atk": context.Value}
}

func (e PieceATKDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetATK(pieceState.piece.ATK() - context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"delta_atk": -context.Value}
}

func (e PieceATKFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) (*model.Board, logic.Summary) {
	pieceState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board, nil
	}
	board = board.Next()
	pieceState.piece.SetATK(context.Value)
	board.UpdatePiece(pieceState.piece)
	return board, logic.Summary{"set_atk": context.Value}
}

package piece_action

import (
	"fmt"
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

func (c PieceSummonActionContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["point"] = pkg.PointToMap(c.Point)
	if c.Piece != nil {
		result["piece_id"] = c.Piece.Id().String()
	}
	return result
}

func (c *PieceSummonActionContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["point"]; ok {
		p, err := pkg.PointFromMap(v)
		if err != nil {
			return fmt.Errorf("point: %w", err)
		}
		c.Point = p
	}
	return nil
}

func (c PieceMoveActionContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["point"] = pkg.PointToMap(c.Point)
	return result
}

func (c *PieceMoveActionContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["point"]; ok {
		p, err := pkg.PointFromMap(v)
		if err != nil {
			return fmt.Errorf("point: %w", err)
		}
		c.Point = p
	}
	return nil
}

func (c PieceAttackActionContext) ToMap() map[string]interface{} {
	result := c.BaseEffectContext.ToMap()
	result["point"] = pkg.PointToMap(c.Point)
	result["value"] = c.Value
	return result
}

func (c *PieceAttackActionContext) FromMap(m map[string]interface{}) error {
	if err := c.BaseEffectContext.FromMap(m); err != nil {
		return err
	}
	if v, ok := m["point"]; ok {
		p, err := pkg.PointFromMap(v)
		if err != nil {
			return fmt.Errorf("point: %w", err)
		}
		c.Point = p
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

func (PieceSummonAction) Tag() logic.EffectActionTag   { return action.ENTITY_SUMMON_ACTION }
func (PieceMoveAction) Tag() logic.EffectActionTag     { return action.ENTITY_MOVE_ACTION }
func (PieceAttackAction) Tag() logic.EffectActionTag   { return action.ENTITY_ATTACK_ACTION }
func (PieceInvasionAction) Tag() logic.EffectActionTag { return action.ENTITY_INVASION_ACTION }

func (e PieceSummonAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(PieceSummonActionState); ok {
		return &PieceSummonActionContext{event.BaseEffectContext{}, s.point, s.piece}, logic.Summary{"point": pkg.PointToMap(s.point)}
	}
	return nil, nil
}

func (e PieceSummonAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	_, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil, nil
	}
	board = board.Next()
	if !board.SetPiece(c.Point, c.Piece) {
		return nil, nil
	}
	return board, logic.Summary{"placed_at": pkg.PointToMap(c.Point)}
}

func (e PieceMoveAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(PieceMoveActionState); ok {
		return &PieceSummonActionContext{event.BaseEffectContext{}, s.to, s.piece}, logic.Summary{"from": pkg.PointToMap(s.from), "to": pkg.PointToMap(s.to)}
	}
	return nil, nil
}

func (e PieceMoveAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil, nil
	}
	board = board.Next()

	if !board.SetPiece(c.Point, s.piece) {
		return nil, nil
	}
	return board, logic.Summary{"to": pkg.PointToMap(c.Point)}
}

func (e PieceAttackAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) (logic.EffectContext, logic.Summary) {
	if s, ok := state.(PieceMoveActionState); ok {
		return &PieceSummonActionContext{event.BaseEffectContext{}, s.to, s.piece}, logic.Summary{"target": pkg.PointToMap(s.to)}
	}
	return nil, nil
}

func (e PieceAttackAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) (*model.Board, logic.Summary) {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil, nil
	}
	board = board.Next()
	s.decreaseHPState.piece = board.Entities()[c.Point.X][c.Point.Y].Copy()
	s.decreaseHPState.value = c.Value
	return board, logic.Summary{"target": pkg.PointToMap(c.Point), "value": c.Value}
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

package entity_action

import (
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/impl/action"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
)

type EntitySummonActionState struct {
	point  pkg.Point
	entity model.Entity
}

type EntityMoveActionState struct {
	from, to pkg.Point
	entity   model.Entity
}

type EntityAttackActionState struct {
	decreaseHPState *EntityOperateActionState
	point           pkg.Point
	entity          model.Entity
	value           int
}

type EntitySummonActionContext struct {
	event.BaseEffectContext
	Point  pkg.Point
	Entity model.Entity
}

type EntityMoveActionContext struct {
	event.BaseEffectContext
	Point pkg.Point
}

type EntityAttackActionContext struct {
	event.BaseEffectContext
	Point pkg.Point
	Value int
}

type EntitySummonAction struct {
	logic.BaseAction[EntitySummonActionState, EntitySummonActionContext]
}

type EntityMoveAction struct {
	logic.BaseAction[EntityMoveActionState, EntityMoveActionContext]
}

type EntityAttackAction struct {
	logic.BaseAction[EntityAttackActionState, EntityAttackActionContext]
}

type EntityInvasionAction struct {
	logic.BaseAction[EntityAttackActionState, EntityAttackActionContext]
}

func (e EntitySummonAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(EntitySummonActionState); ok {
		return &EntitySummonActionContext{event.BaseEffectContext{}, s.point, s.entity}
	}
	return nil
}

func (e EntitySummonAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	_, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()
	if !board.SetEntity(c.Point, c.Entity) {
		return nil
	}
	return board
}

func (e EntityMoveAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(EntityMoveActionState); ok {
		return &EntitySummonActionContext{event.BaseEffectContext{}, s.to, s.entity}
	}
	return nil
}

func (e EntityMoveAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()

	if !board.SetEntity(c.Point, s.entity) {
		return nil
	}
	return board
}

func (e EntityAttackAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(EntityMoveActionState); ok {
		return &EntitySummonActionContext{event.BaseEffectContext{}, s.to, s.entity}
	}
	return nil
}

func (e EntityAttackAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()
	s.decreaseHPState.entity = board.Entities()[c.Point.X][c.Point.Y].Copy()
	s.decreaseHPState.value = c.Value
	return board
}

func (e EntityAttackAction) SubEffects(state interface{}) []*logic.EffectEvent {
	s, ok := state.(EntityAttackActionState)
	if !ok {
		return nil
	}
	result := make([]*logic.EffectEvent, 1)
	result[0] = logic.NewEffectEvent(action.ENTITY_HP_DECREASE_ACTION, s.decreaseHPState)
	return result
}

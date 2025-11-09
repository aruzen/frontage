package entity_action

import (
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/logic"
	"frontage/pkg/model"
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
	point  pkg.Point
	entity model.Entity
	value  int
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

type EntityInvationAction struct {
	logic.BaseAction[EntityAttackActionState, EntityAttackActionContext]
}

func (e EntitySummonAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(EntitySummonActionState); ok {
		return EntitySummonActionContext{event.BaseEffectContext{}, s.point, s.entity}
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

func (e EntityMoveAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(EntityMoveActionState); ok {
		return EntitySummonActionContext{event.BaseEffectContext{}, s.to, s.entity}
	}
	return nil
}

func (e EntityMoveAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()
	moving := board.Entities()[s.from.X][s.from.Y]
	if moving == nil {
		return nil
	}
	if !board.RemoveEntity(s.from) {
		return nil
	}
	if !board.SetEntity(c.Point, moving) {
		return nil
	}
	return board
}

func (e EntityAttackAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if s, ok := state.(EntityMoveActionState); ok {
		return EntitySummonActionContext{event.BaseEffectContext{}, s.to, s.entity}
	}
	return nil
}

func (e EntityAttackAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) *model.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()

	return board
}

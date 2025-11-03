package entity

import (
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/entity"
	"frontage/pkg/event/action"
)

type EntitySummonActionState struct {
	point  pkg.Point
	entity entity.Entity
}

type EntityMoveActionState struct {
	from, to pkg.Point
	entity   entity.Entity
}

type EntityAttackActionState struct {
	point  pkg.Point
	entity entity.Entity
	value  int
}

type EntitySummonActionContext struct {
	event.BaseEffectContext
	Point  pkg.Point
	Entity entity.Entity
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
	event.BaseAction[EntitySummonActionState, EntitySummonActionContext]
}

type EntityMoveAction struct {
	event.BaseAction[EntityMoveActionState, EntityMoveActionContext]
}

type EntityAttackAction struct {
	event.BaseAction[EntityAttackActionState, EntityAttackActionContext]
}

type EntityInvationAction struct {
	event.BaseAction[EntityAttackActionState, EntityAttackActionContext]
}

func (e EntitySummonAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if s, ok := state.(EntitySummonActionState); ok {
		return EntitySummonActionContext{event.BaseEffectContext{}, s.point, s.entity}
	}
	return nil
}

func (e EntitySummonAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) *pkg.Board {
	_, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()
	board.Entities()[c.Point.X][c.Point.Y] = c.Entity
	return board
}

func (e EntityMoveAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if s, ok := state.(EntityMoveActionState); ok {
		return EntitySummonActionContext{event.BaseEffectContext{}, s.to, s.entity}
	}
	return nil
}

func (e EntityMoveAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) *pkg.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()
	board.Entities()[c.Point.X][c.Point.Y] = board.Entities()[s.from.X][s.from.Y]
	board.Entities()[s.from.X][s.from.Y] = nil
	return board
}

func (e EntityAttackAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if s, ok := state.(EntityMoveActionState); ok {
		return EntitySummonActionContext{event.BaseEffectContext{}, s.to, s.entity}
	}
	return nil
}

func (e EntityAttackAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) *pkg.Board {
	s, c, ok := e.CastStateContext(state, context)
	if !ok {
		return nil
	}
	board = board.Next()

	return board
}

package entity_action

import (
	"frontage/internal/event"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"log/slog"
)

type EntityOperateActionState struct {
	entity model.MutableEntity
	value  int
}

type EntityHPContext struct {
	event.BaseEffectContext
	Value int
}

type EntityMPContext struct {
	event.BaseEffectContext
	Value int
}

type EntityATKContext struct {
	event.BaseEffectContext
	Value int
}

type baseEntityHPAction struct {
	logic.BaseAction[EntityOperateActionState, EntityHPContext]
}

type baseEntityMPAction struct {
	logic.BaseAction[EntityOperateActionState, EntityMPContext]
}

type baseEntityATKAction struct {
	logic.BaseAction[EntityOperateActionState, EntityATKContext]
}

type EntityHPIncreaseAction struct{ baseEntityHPAction }

type EntityHPDecreaseAction struct{ baseEntityHPAction }

type EntityHPFixAction struct{ baseEntityHPAction }

type EntityMPIncreaseAction struct{ baseEntityMPAction }

type EntityMPDecreaseAction struct{ baseEntityMPAction }

type EntityMPFixAction struct{ baseEntityMPAction }

type EntityATKIncreaseAction struct{ baseEntityATKAction }

type EntityATKDecreaseAction struct{ baseEntityATKAction }

type EntityATKFixAction struct{ baseEntityATKAction }

func (e EntityOperateActionState) Entity() model.MutableEntity {
	return e.entity
}

func (e EntityOperateActionState) Value() int {
	return e.value
}

func (b baseEntityHPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return &EntityHPContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityHPIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetHP(entityState.entity.HP() + context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (e EntityHPDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetHP(entityState.entity.HP() - context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (e EntityHPFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetHP(context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (b baseEntityMPAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return &EntityMPContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityMPIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetMP(entityState.entity.MP() + context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (e EntityMPDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetMP(entityState.entity.MP() - context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (e EntityMPFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetMP(context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (b baseEntityATKAction) Act(state interface{}, beforeAction logic.EffectAction, beforeContext logic.EffectContext) logic.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return &EntityATKContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityATKIncreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetATK(entityState.entity.ATK() + context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (e EntityATKDecreaseAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetATK(entityState.entity.ATK() - context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

func (e EntityATKFixAction) Solve(board *model.Board, state interface{}, c logic.EffectContext) *model.Board {
	entityState, context, ok := e.CastStateContext(state, c)
	if !ok {
		slog.Warn("CastStateContext failed.")
		return board
	}
	board = board.Next()
	entityState.entity.SetATK(context.Value)
	board.UpdateEntity(entityState.entity)
	return board
}

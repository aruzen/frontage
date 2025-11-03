package entity

import (
	"frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/entity"
	"frontage/pkg/event/action"
	"log/slog"
)

type EntityOperateActionState struct {
	entity entity.MutableEntity
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
	event.BaseAction[EntityOperateActionState, EntityHPContext]
}

type baseEntityMPAction struct {
	event.BaseAction[EntityOperateActionState, EntityMPContext]
}

type baseEntityATKAction struct {
	event.BaseAction[EntityOperateActionState, EntityATKContext]
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

func (e EntityOperateActionState) Entity() entity.MutableEntity {
	return e.entity
}

func (e EntityOperateActionState) Value() int {
	return e.value
}

func (b baseEntityHPAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return EntityHPContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityHPIncreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, hpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetHP(entityState.entity.HP() + hpContext.Value)
	}
}

func (e EntityHPDecreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, hpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetHP(entityState.entity.HP() - hpContext.Value)
	}
}

func (e EntityHPFixAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, hpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetHP(hpContext.Value)
	}
}

func (b baseEntityMPAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return EntityMPContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityMPIncreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, mpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetMP(entityState.entity.MP() + mpContext.Value)
	}
}

func (e EntityMPDecreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, mpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetMP(entityState.entity.MP() - mpContext.Value)
	}
}

func (e EntityMPFixAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, mpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetMP(mpContext.Value)
	}
}

func (b baseEntityATKAction) Act(state interface{}, beforeContext action.EffectContext) action.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return EntityATKContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityATKIncreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, atkContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetATK(entityState.entity.ATK() + atkContext.Value)
	}
}

func (e EntityATKDecreaseAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, atkContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetATK(entityState.entity.ATK() - atkContext.Value)
	}
}

func (e EntityATKFixAction) Solve(board *pkg.Board, state interface{}, context action.EffectContext) {
	entityState, atkContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetATK(atkContext.Value)
	}
}

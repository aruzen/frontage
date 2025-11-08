package entity_action

import (
	"frontage/internal/event"
	"frontage/pkg/logic"
	"frontage/pkg/model"
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

func (b baseEntityHPAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return EntityHPContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityHPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, hpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetHP(entityState.entity.HP() + hpContext.Value)
	}
}

func (e EntityHPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, hpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetHP(entityState.entity.HP() - hpContext.Value)
	}
}

func (e EntityHPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, hpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetHP(hpContext.Value)
	}
}

func (b baseEntityMPAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return EntityMPContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityMPIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, mpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetMP(entityState.entity.MP() + mpContext.Value)
	}
}

func (e EntityMPDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, mpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetMP(entityState.entity.MP() - mpContext.Value)
	}
}

func (e EntityMPFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, mpContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetMP(mpContext.Value)
	}
}

func (b baseEntityATKAction) Act(state interface{}, beforeContext logic.EffectContext) logic.EffectContext {
	if entityState, ok := state.(EntityOperateActionState); ok {
		return EntityATKContext{event.NewBaseEffectContext(), entityState.value}
	}
	slog.Warn("State was not EntityOperateActionState.")
	return nil
}

func (e EntityATKIncreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, atkContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetATK(entityState.entity.ATK() + atkContext.Value)
	}
}

func (e EntityATKDecreaseAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, atkContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetATK(entityState.entity.ATK() - atkContext.Value)
	}
}

func (e EntityATKFixAction) Solve(board *model.Board, state interface{}, context logic.EffectContext) {
	entityState, atkContext, ok := e.CastStateContext(state, context)
	if ok {
		entityState.entity.SetATK(atkContext.Value)
	}
}

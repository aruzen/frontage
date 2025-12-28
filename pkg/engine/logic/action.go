package logic

import (
	event_base "frontage/internal/event"
	"frontage/pkg"
	"frontage/pkg/engine/model"
	"reflect"
)

type EffectActionTag pkg.LocalizeTag
type ModifyActionTag pkg.LocalizeTag

type Action interface {
	pkg.Localized
	WantState() reflect.Type
	WantContext() reflect.Type
}

type EffectContext interface {
	IsCanceled() bool
	Cancel()
	ToMap() map[string]interface{}
	FromMap(map[string]interface{}) error
}

type ActionState interface {
	ToMap() map[string]interface{}
	FromMap(*model.Board, map[string]interface{}) error
}

type Summary map[string]interface{}

type SummaryType string

const (
	SUMMARY_TYPE_ACT    SummaryType = "act"
	SUMMARY_TYPE_MODIFY SummaryType = "modify"
	SUMMARY_TYPE_SOLVE  SummaryType = "solve"
)

type ActionSummary struct {
	Action Action
	Type   SummaryType
	Data   Summary
}

type ActionResult struct {
	Action     Action
	State      ActionState
	Context    EffectContext
	SummaryIdx int
}

type EffectAction interface {
	Action
	Tag() EffectActionTag
	Act(state ActionState, beforeAction EffectAction, beforeContext EffectContext) (EffectContext, Summary)
	Solve(board *model.Board, state ActionState, context EffectContext) (*model.Board, Summary)
	MakeContext(data map[string]interface{}) EffectContext
	MakeState(board *model.Board, data map[string]interface{}) ActionState
}

type ModifyAction interface {
	Action
	Tag() ModifyActionTag
	Modify(state ActionState, context EffectContext) (EffectContext, Summary)
}

type MultiEffectAction interface {
	Action
	SubEffects(state ActionState) []*EffectEvent
}

type BaseAction[StateType any, ContextType any] struct {
	event_base.BaseWantState[StateType]
	event_base.BaseWantContext[ContextType]
}

// LocalizeTag はデフォルトで空文字を返す。具体的なタグを持つ Action は個別に上書きすること。
func (a BaseAction[StateType, ContextType]) LocalizeTag() pkg.LocalizeTag {
	return ""
}

func (a *BaseAction[StateType, ContextType]) CastStateContext(state ActionState, context EffectContext) (*StateType, *ContextType, bool) {
	if state == nil || context == nil {
		return nil, nil, false
	}
	var castedState *StateType
	switch s := any(state).(type) {
	case StateType:
		castedState = &s
	case *StateType:
		castedState = s
	default:
		return nil, nil, false
	}
	var castedContext *ContextType
	switch c := any(context).(type) {
	case ContextType:
		castedContext = &c
	case *ContextType:
		castedContext = c
	default:
		return nil, nil, false
	}
	return castedState, castedContext, true
}

func (a BaseAction[StateType, ContextType]) MakeContext(data map[string]interface{}) EffectContext {
	var ctx ContextType
	if v, ok := any(&ctx).(EffectContext); ok {
		if err := v.FromMap(data); err != nil {
			return nil
		}
		return v
	}
	if v, ok := any(ctx).(EffectContext); ok {
		if err := v.FromMap(data); err != nil {
			return nil
		}
		return v
	}
	return nil
}

func (a BaseAction[StateType, ContextType]) MakeState(board *model.Board, data map[string]interface{}) ActionState {
	var st StateType
	if v, ok := any(&st).(ActionState); ok {
		if err := v.FromMap(board, data); err != nil {
			return nil
		}
		return v
	}
	if v, ok := any(st).(ActionState); ok {
		if err := v.FromMap(board, data); err != nil {
			return nil
		}
		return v
	}
	return nil
}

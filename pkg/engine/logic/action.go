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

type ActionResult struct {
	Action  Action
	Context EffectContext
}

type EffectAction interface {
	Action
	Tag() EffectActionTag
	Act(state interface{}, beforeAction EffectAction, beforeContext EffectContext) EffectContext
	Solve(board *model.Board, state interface{}, context EffectContext) *model.Board
}

type ModifyAction interface {
	Action
	Tag() ModifyActionTag
	Modify(state interface{}, context EffectContext) EffectContext
}

type MultiEffectAction interface {
	Action
	SubEffects(state interface{}) []*EffectEvent
}

type BaseAction[StateType any, ContextType any] struct {
	event_base.BaseWantState[StateType]
	event_base.BaseWantContext[ContextType]
}

// ItemTag はデフォルトで空文字を返す。具体的なタグを持つ Action は個別に上書きすること。
func (a *BaseAction[StateType, ContextType]) ItemTag() pkg.LocalizeTag {
	return ""
}

func (a *BaseAction[StateType, ContextType]) Children() []EffectAction {
	return nil
}

func (a *BaseAction[StateType, ContextType]) CastStateContext(state interface{}, context interface{}) (*StateType, *ContextType, bool) {
	castedState, ok := state.(StateType)
	if !ok {
		return nil, nil, false
	}
	castedContext, ok := context.(ContextType)
	if !ok {
		return nil, nil, false
	}
	return &castedState, &castedContext, true
}

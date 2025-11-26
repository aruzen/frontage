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
	Action  Action
	Context EffectContext
}

type EffectAction interface {
	Action
	Tag() EffectActionTag
	Act(state interface{}, beforeAction EffectAction, beforeContext EffectContext) (EffectContext, Summary)
	Solve(board *model.Board, state interface{}, context EffectContext) (*model.Board, Summary)
}

type ModifyAction interface {
	Action
	Tag() ModifyActionTag
	Modify(state interface{}, context EffectContext) (EffectContext, Summary)
}

type MultiEffectAction interface {
	Action
	SubEffects(state interface{}) []*EffectEvent
}

type BaseAction[StateType any, ContextType any] struct {
	event_base.BaseWantState[StateType]
	event_base.BaseWantContext[ContextType]
}

// LocalizeTag はデフォルトで空文字を返す。具体的なタグを持つ Action は個別に上書きすること。
func (a BaseAction[StateType, ContextType]) LocalizeTag() pkg.LocalizeTag {
	return ""
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

package event

import (
	"frontage/pkg/event/action"
	"reflect"
)

type BaseEffectContext struct {
	canceled bool
}

func NewBaseEffectContext() BaseEffectContext {
	return BaseEffectContext{canceled: false}
}

func (b BaseEffectContext) IsCanceled() bool {
	return b.canceled
}

func (b BaseEffectContext) Cancel() {
	b.canceled = true
}

type BaseWantState[StateType interface{}] struct {
}

func (a *BaseWantState[StateType]) WantState() reflect.Type {
	return reflect.TypeFor[StateType]()
}

type BaseWantContext[ContextType interface{}] struct {
}

func (a *BaseWantContext[ContextType]) WantContext() reflect.Type {
	return reflect.TypeFor[ContextType]()
}

type BaseAction[StateType interface{}, ContextType interface{}] struct {
	BaseWantState[StateType]
	BaseWantContext[ContextType]
}

func (a *BaseAction[StateType, ContextType]) Children() []action.EffectAction {
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

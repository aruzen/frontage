package event

import (
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

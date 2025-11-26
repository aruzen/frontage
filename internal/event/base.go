package event

import (
	"fmt"
	"reflect"
)

type BaseEffectContext struct {
	canceled bool
}

func NewBaseEffectContext() BaseEffectContext {
	return BaseEffectContext{canceled: false}
}

func (b *BaseEffectContext) IsCanceled() bool {
	return b.canceled
}

func (b *BaseEffectContext) Cancel() {
	b.canceled = true
}

func (b *BaseEffectContext) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"canceled": b.canceled,
	}
}

func (b *BaseEffectContext) FromMap(m map[string]interface{}) error {
	if m == nil {
		return nil
	}
	if v, ok := m["canceled"]; ok {
		switch val := v.(type) {
		case bool:
			b.canceled = val
		case nil:
			b.canceled = false
		default:
			return fmt.Errorf("canceled expects bool, got %T", v)
		}
	}
	return nil
}

type BaseWantState[StateType any] struct {
}

func (a BaseWantState[StateType]) WantState() reflect.Type {
	return reflect.TypeFor[StateType]()
}

type BaseWantContext[ContextType any] struct {
}

func (a BaseWantContext[ContextType]) WantContext() reflect.Type {
	return reflect.TypeFor[ContextType]()
}

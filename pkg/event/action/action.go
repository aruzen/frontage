package action

import (
	"frontage/pkg"
	"frontage/pkg/event"
	"reflect"
)

type Action interface {
	WantState() reflect.Type
	WantContext() reflect.Type
}

type EffectContext interface {
	IsCanceled() bool
	Cancel()
}

type EffectAction interface {
	Action
	Act(state interface{}, beforeAction EffectAction, beforeContext EffectContext) EffectContext
	Solve(board *pkg.Board, state interface{}, context EffectContext) *pkg.Board
}

type ModifyAction interface {
	Action
	Modify(state interface{}, context EffectContext) EffectContext
}

type MultiEffectAction interface {
	SubEffects(state interface{}) []event.EffectEvent
}

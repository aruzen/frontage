package action

import (
	"frontage/pkg"
	"reflect"
)

type Action interface {
	Children() []EffectAction
	WantState() reflect.Type
	WantContext() reflect.Type
}

type EffectContext interface {
	IsCanceled() bool
	Cancel()
}

type EffectAction interface {
	Action
	Act(state interface{}, beforeContext EffectContext) EffectContext
	Solve(board *pkg.Board, state interface{}, context EffectContext) *pkg.Board
}

type ModifyAction interface {
	Action
	Modify(state interface{}, context EffectContext) EffectContext
}

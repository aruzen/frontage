package logic

import (
	"testing"

	"frontage/internal/event"
	"frontage/pkg/engine/model"
)

type abortProbeState struct{}

func (*abortProbeState) ToMap() map[string]interface{}                      { return nil }
func (*abortProbeState) FromMap(*model.Board, map[string]interface{}) error { return nil }

type abortProbeAction struct {
	BaseAction[abortProbeState, event.BaseEffectContext]
	nilAct, nilSolve bool
	solves           int
}

func (*abortProbeAction) Tag() EffectActionTag { return "abort_probe" }
func (a *abortProbeAction) Act(ActionState, EffectAction, EffectContext) (EffectContext, Summary) {
	if a.nilAct {
		return nil, nil
	}
	context := event.NewBaseEffectContext()
	return &context, nil
}
func (a *abortProbeAction) Solve(board *model.Board, _ ActionState, _ EffectContext) (*model.Board, Summary) {
	a.solves++
	if a.nilSolve {
		return nil, nil
	}
	return board, nil
}

type abortProbeModifier struct {
	BaseAction[abortProbeState, event.BaseEffectContext]
	nilResult bool
	calls     int
}

func (*abortProbeModifier) Tag() ModifyActionTag { return "abort_probe_modify" }
func (a *abortProbeModifier) Modify(_ ActionState, context EffectContext) (EffectContext, Summary) {
	a.calls++
	if a.nilResult {
		return nil, nil
	}
	if context == nil {
		fresh := event.NewBaseEffectContext()
		context = &fresh
	}
	return context, nil
}

func TestNilSolveAbortsDependentEffects(t *testing.T) {
	for _, fail := range []bool{true, false} {
		board := &model.Board{}
		es := NewEventSystem(board)
		parent := &abortProbeAction{nilSolve: fail}
		child := &abortProbeAction{}
		root := NewEffectEvent(parent, &abortProbeState{})
		root.branch = []*EffectEvent{NewEffectEvent(child, &abortProbeState{})}
		root.Resolve(es, nil, nil)
		if es.Board != board {
			t.Fatal("nil Solve replaced the valid board")
		}
		if fail {
			if child.solves != 0 || len(es.AppliedEffects) != 0 {
				t.Fatal("failed effect ran a dependent effect or was recorded as applied")
			}
		} else if child.solves != 1 || len(es.AppliedEffects) != 2 {
			t.Fatal("successful effect no longer resolves its children")
		}
	}
}

func TestNilActCannotBeRevivedByModifier(t *testing.T) {
	action := &abortProbeAction{nilAct: true}
	modifier := &abortProbeModifier{}
	root := NewEffectEvent(action, &abortProbeState{})
	root.modifier = NewModifyEvent(modifier, &abortProbeState{})
	es := NewEventSystem(&model.Board{})
	root.Resolve(es, nil, nil)
	if modifier.calls != 0 || action.solves != 0 || len(es.AppliedEffects) != 0 {
		t.Fatal("nil Act continued to a modifier or Solve")
	}
}

func TestNilModifyAbortsFollowingModifiers(t *testing.T) {
	action := &abortProbeAction{}
	first := &abortProbeModifier{nilResult: true}
	second := &abortProbeModifier{}
	root := NewEffectEvent(action, &abortProbeState{})
	root.modifier = NewModifyEvent(first, &abortProbeState{})
	root.modifier.modifier = NewModifyEvent(second, &abortProbeState{})
	es := NewEventSystem(&model.Board{})
	root.Resolve(es, nil, nil)
	if first.calls != 1 || second.calls != 0 || action.solves != 0 || len(es.AppliedEffects) != 0 {
		t.Fatal("nil Modify continued to subsequent processing")
	}
}

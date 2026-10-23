package logic

import (
	"testing"

	"frontage/internal/event"
	"frontage/pkg/engine/model"
)

type resultState struct{}

func (*resultState) ToMap() map[string]interface{}                      { return nil }
func (*resultState) FromMap(*model.Board, map[string]interface{}) error { return nil }

type resultAction struct {
	BaseAction[resultState, event.BaseEffectContext]
	result EffectContext
	solves int
}

func (*resultAction) Tag() EffectActionTag { return "test/result" }
func (a *resultAction) Act(ActionState, EffectAction, EffectContext) (EffectContext, Summary) {
	return a.result, nil
}
func (a *resultAction) Solve(board *model.Board, _ ActionState, _ EffectContext) (*model.Board, Summary) {
	a.solves++
	return board, nil
}

type resultModifier struct {
	BaseAction[resultState, event.BaseEffectContext]
	result EffectContext
	calls  int
}

func (*resultModifier) Tag() ModifyActionTag { return "test/result_modify" }
func (m *resultModifier) Modify(ActionState, EffectContext) (EffectContext, Summary) {
	m.calls++
	return m.result, nil
}

func newResultContext() EffectContext {
	context := event.NewBaseEffectContext()
	return &context
}

func resolveResult(t *testing.T, root *EffectEvent, system *EventSystem) {
	t.Helper()
	defer func() {
		if value := recover(); value != nil {
			t.Errorf("result context caused panic: %v", value)
		}
	}()
	root.Resolve(system, nil, nil)
}

func TestNilResultContextStopsResolution(t *testing.T) {
	var typedNil *event.BaseEffectContext
	for _, tc := range []struct {
		name    string
		context EffectContext
		valid   bool
	}{
		{"nil", nil, false},
		{"typed_nil", typedNil, false},
		{"valid", newResultContext(), true},
	} {
		for _, phase := range []string{"act", "modify"} {
			t.Run(tc.name+"/"+phase, func(t *testing.T) {
				parent := &resultAction{result: newResultContext()}
				child := &resultAction{result: newResultContext()}
				first := &resultModifier{result: newResultContext()}
				second := &resultModifier{result: newResultContext()}
				root := NewEffectEvent(parent, &resultState{})
				root.modifier = NewModifyEvent(first, &resultState{})
				root.modifier.branch = []*EffectEvent{NewEffectEvent(child, &resultState{})}
				root.modifier.modifier = NewModifyEvent(second, &resultState{})
				if phase == "act" {
					parent.result = tc.context
				} else {
					first.result = tc.context
				}
				board := &model.Board{}
				system := NewEventSystem(board)
				resolveResult(t, root, system)
				if system.Board != board {
					t.Error("failed context changed the live board")
				}
				if tc.valid {
					if parent.solves != 1 || child.solves != 1 || second.calls != 1 || len(system.AppliedEffects) != 2 {
						t.Error("valid context no longer completes resolution")
					}
					return
				}
				if parent.solves != 0 || child.solves != 0 || second.calls != 0 || len(system.AppliedEffects) != 0 {
					t.Error("nil result allowed dependent processing or was recorded as applied")
				}
				if phase == "act" && first.calls != 0 {
					t.Error("modifier ran after a failed Act")
				}
			})
		}
	}
}

func TestTypedNilResultIsNotDereferenced(t *testing.T) {
	var context *event.BaseEffectContext
	for _, phase := range []string{"act", "modify"} {
		t.Run(phase, func(t *testing.T) {
			parent := &resultAction{result: context}
			root := NewEffectEvent(parent, &resultState{})
			if phase == "modify" {
				parent.result = newResultContext()
				root.modifier = NewModifyEvent(&resultModifier{result: context}, &resultState{})
			}
			board := &model.Board{}
			system := NewEventSystem(board)
			resolveResult(t, root, system)
			if parent.solves != 0 || system.Board != board || len(system.AppliedEffects) != 0 {
				t.Error("typed nil result did not abort resolution")
			}
		})
	}
}

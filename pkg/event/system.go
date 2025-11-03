package event

import (
	"frontage/pkg"
	"frontage/pkg/event/action"
	"log/slog"
	"reflect"
)

type EventSystem struct {
	Board         *pkg.Board
	Trigger       EffectEvent
	ResultEffects []struct {
		Action  action.Action
		Context action.EffectContext
	}
}

type baseEvent struct {
	branch   []*EffectEvent
	modifier *ModifyEvent
	state    interface{}
}

type EffectEvent struct {
	baseEvent
	action action.EffectAction
}

type ModifyEvent struct {
	baseEvent
	action action.ModifyAction
}

/*
func NewEvent(parent *Event, action Action, context interface{}) *Event {
	return &Event{parent, action, context}
}
*/

func (e *baseEvent) Branch() []*EffectEvent {
	return e.branch
}

func (e *baseEvent) Modifier() *ModifyEvent {
	return e.modifier
}

func (e *baseEvent) State() interface{} {
	return e.state
}

func (e *EffectEvent) Action() action.EffectAction {
	return e.action
}

func (e *ModifyEvent) Action() action.ModifyAction {
	return e.action
}

func IntegrityCheck(a action.Action, state interface{}) bool {
	if a == nil {
		return false
	}
	if reflect.TypeOf(state) != a.WantState() {
		return false
	}
	return true
}

func (es *EventSystem) Resolve() {

}

func (e *EffectEvent) Resolve(es *EventSystem, effect action.EffectAction, beforeContext action.EffectContext) {
	if IntegrityCheck(e.Action(), e.State()) {
		slog.Warn("warning: IntegrityCheck failed.")
		return
	}
	context := e.Action().Act(e.State(), beforeContext)
	if e.Modifier() != nil && IntegrityCheck(e.Modifier().Action(), e.Modifier().State()) {
		context = e.Modifier().Resolve(es, e.Action(), context)
	}
	for _, branch := range e.Branch() {
		if IntegrityCheck(branch.action, branch.state) {
			slog.Warn("warning: IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, effect, context)
	}
	if !context.IsCanceled() {
		es.ResultEffects = append(es.ResultEffects, struct {
			Action  action.Action
			Context action.EffectContext
		}{Action: e.action, Context: context})
	}
}

func (e *ModifyEvent) Resolve(es *EventSystem, effect action.EffectAction, context action.EffectContext) action.EffectContext {
	if IntegrityCheck(e.Action(), e.State()) {
		slog.Warn("warning: IntegrityCheck failed.")
		return nil
	}
	context = e.Action().Modify(e.State(), context)
	for _, branch := range e.Branch() {
		if IntegrityCheck(branch.action, branch.state) {
			slog.Warn("warning: IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, effect, context)
	}
	return context
}

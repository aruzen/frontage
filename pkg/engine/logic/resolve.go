package logic

import (
	"frontage/pkg/engine/model"
	"log/slog"
)

func (es *EventSystem) Resolve() {
	es.trigger.Resolve(es, nil, nil)
}

func (es *EventSystem) ResolveTest() (*model.Board, []ActionResult) {
	sandbox := &EventSystem{
		Board:          es.Board.Sandbox(),
		trigger:        es.trigger,
		active:         es.active,
		AppliedEffects: make([]ActionResult, 0),
	}
	sandbox.Resolve()
	return es.Board, es.AppliedEffects
}

func (e *EffectEvent) Resolve(es *EventSystem, beforeEffect EffectAction, beforeContext EffectContext) {
	if !IntegrityCheck(e.Action(), e.State()) {
		slog.Warn("IntegrityCheck failed.")
		return
	}
	context := e.action.Act(e.State(), beforeEffect, beforeContext)
	if e.modifier != nil && IntegrityCheck(e.modifier.action, *e.modifier.state) {
		context = e.modifier.Resolve(es, e.action, context)
	}

	if context.IsCanceled() {
		return
	}

	es.Board = e.action.Solve(es.Board, e.State(), context)

	for _, branch := range e.branch {
		if !IntegrityCheck(branch.action, *branch.state) {
			slog.Warn("IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, e.action, context)
	}

	es.AppliedEffects = append(es.AppliedEffects, struct {
		Action  Action
		Context EffectContext
	}{Action: e.action, Context: context})
}

func (e *ModifyEvent) Resolve(es *EventSystem, effect EffectAction, context EffectContext) EffectContext {
	if !IntegrityCheck(e.action, *e.state) {
		slog.Warn("IntegrityCheck failed.")
		return nil
	}
	context = e.action.Modify(*e.state, context)
	for _, branch := range e.branch {
		if !IntegrityCheck(branch.action, *branch.state) {
			slog.Warn("IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, effect, context)
	}
	return context
}

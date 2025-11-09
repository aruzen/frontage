package logic

import "log/slog"

func (es *EventSystem) Resolve() {
	es.Trigger.Resolve(es, nil, nil)
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
	for _, branch := range e.branch {
		if !IntegrityCheck(branch.action, *branch.state) {
			slog.Warn("IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, e.action, context)
	}
	if !context.IsCanceled() {
		es.ResultEffects = append(es.ResultEffects, struct {
			Action  Action
			Context EffectContext
		}{Action: e.action, Context: context})
	}
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

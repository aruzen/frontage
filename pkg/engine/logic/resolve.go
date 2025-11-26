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
	var summaries *[]ActionSummary
	var tmpSummary Summary
	*summaries = make([]ActionSummary, 0, 3)
	es.Summaries = append(es.Summaries, summaries)

	context, tmpSummary := e.action.Act(e.State(), beforeEffect, beforeContext)
	appendSummary(summaries, e.action, SUMMARY_TYPE_ACT, tmpSummary)

	if e.modifier != nil && IntegrityCheck(e.modifier.action, *e.modifier.state) {
		context = e.modifier.Resolve(es, e.action, context, summaries)
	}

	if context == nil || context.IsCanceled() {
		return
	}

	es.Board, tmpSummary = e.action.Solve(es.Board, e.State(), context)
	appendSummary(summaries, e.action, SUMMARY_TYPE_SOLVE, tmpSummary)

	for _, branch := range e.branch {
		if !IntegrityCheck(branch.action, *branch.state) {
			slog.Warn("IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, e.action, context)
	}

	es.AppliedEffects = append(es.AppliedEffects, ActionResult{
		Action:  e.action,
		Context: context,
	})
}

func (e *ModifyEvent) Resolve(es *EventSystem, effect EffectAction, context EffectContext, summaries *[]ActionSummary) EffectContext {
	if !IntegrityCheck(e.action, *e.state) {
		slog.Warn("IntegrityCheck failed.")
		return nil
	}
	context, tmpSummary := e.action.Modify(*e.state, context)
	appendSummary(summaries, e.action, SUMMARY_TYPE_MODIFY, tmpSummary)

	for _, branch := range e.branch {
		if !IntegrityCheck(branch.action, *branch.state) {
			slog.Warn("IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, effect, context)
	}

	if e.modifier != nil && IntegrityCheck(e.modifier.action, *e.modifier.state) {
		context = e.modifier.Resolve(es, effect, context, summaries)
	}

	return context
}

func appendSummary(summaries *[]ActionSummary, action Action, st SummaryType, summary Summary) {
	*summaries = append(*summaries, ActionSummary{Action: action, Type: st, Data: summary})
}

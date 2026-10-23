package logic

import (
	"fmt"
	"frontage/pkg/engine/model"
	"log/slog"
	"reflect"
)

func (es *EventSystem) Resolve() {
	if es == nil || es.trigger == nil || es.Board == nil {
		return
	}
	es.trigger.Resolve(es, nil, nil)
}

func (es *EventSystem) ResolveTest() (*model.Board, []ActionResult) {
	if es == nil || es.trigger == nil || es.Board == nil {
		return nil, nil
	}
	sandbox := &EventSystem{
		Board:          es.Board.Sandbox(),
		AppliedEffects: make([]ActionResult, 0),
	}
	trigger, err := cloneEventTree(es.trigger, sandbox.Board)
	if err != nil {
		slog.Warn("Failed to clone event for simulation", "err", err)
		return nil, nil
	}
	sandbox.trigger = trigger.(*EffectEvent)
	sandbox.Resolve()
	return sandbox.Board, sandbox.AppliedEffects
}

func cloneEventTree(root Event, board *model.Board) (Event, error) {
	states := make(map[ActionState]ActionState)
	var cloneState func(ActionState) (ActionState, error)
	cloneState = func(state ActionState) (ActionState, error) {
		if state == nil || !reflect.TypeOf(state).Comparable() {
			return nil, fmt.Errorf("cannot clone state %T", state)
		}
		if copied, ok := states[state]; ok {
			return copied, nil
		}
		var copied ActionState
		if cloner, ok := state.(ActionStateCloner); ok {
			var err error
			copied, err = cloner.CloneState(board, cloneState)
			if err != nil {
				return nil, err
			}
		} else {
			typ := reflect.TypeOf(state)
			if typ.Kind() == reflect.Pointer {
				typ = typ.Elem()
			}
			var ok bool
			copied, ok = reflect.New(typ).Interface().(ActionState)
			if !ok {
				return nil, fmt.Errorf("cannot allocate state %T", state)
			}
			if err := copied.FromMap(board, state.ToMap()); err != nil {
				return nil, err
			}
		}
		states[state] = copied
		return copied, nil
	}
	var cloneEvent func(Event, Event) (Event, error)
	cloneEvent = func(event, parent Event) (Event, error) {
		state, err := cloneState(event.State())
		if err != nil {
			return nil, err
		}
		copied := NewEvent(event.Action(), state)
		if copied == nil {
			return nil, fmt.Errorf("cannot clone action %T", event.Action())
		}
		copied.base().parent = parent
		if modifier := event.base().modifier; modifier != nil {
			cloned, err := cloneEvent(modifier, copied)
			if err != nil {
				return nil, err
			}
			copied.base().modifier = cloned.(*ModifyEvent)
		}
		for _, branch := range event.base().branch {
			cloned, err := cloneEvent(branch, copied)
			if err != nil {
				return nil, err
			}
			copied.base().branch = append(copied.base().branch, cloned.(*EffectEvent))
		}
		return copied, nil
	}
	return cloneEvent(root, nil)
}

func (e *EffectEvent) Resolve(es *EventSystem, beforeEffect EffectAction, beforeContext EffectContext) {
	if !IntegrityCheck(e.Action(), e.State()) {
		slog.Warn("IntegrityCheck failed.")
		return
	}
	summaries := make([]ActionSummary, 0, 3)
	var tmpSummary Summary
	summaryIdx := len(es.Summaries)
	es.Summaries = append(es.Summaries, &summaries)

	context, tmpSummary := e.action.Act(e.State(), beforeEffect, beforeContext)
	appendSummary(&summaries, e.action, SUMMARY_TYPE_ACT, tmpSummary)
	if isNil(context) {
		return
	}

	if e.modifier != nil && IntegrityCheck(e.modifier.action, e.modifier.state) {
		context = e.modifier.Resolve(es, e.action, context, &summaries)
	}

	if isNil(context) || context.IsCanceled() {
		return
	}

	es.Board, tmpSummary = e.action.Solve(es.Board, e.State(), context)
	appendSummary(&summaries, e.action, SUMMARY_TYPE_SOLVE, tmpSummary)

	for _, branch := range e.branch {
		if !IntegrityCheck(branch.action, branch.state) {
			slog.Warn("IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, e.action, context)
	}

	es.AppliedEffects = append(es.AppliedEffects, ActionResult{
		Action:     e.action,
		State:      e.State(),
		Context:    context,
		SummaryIdx: summaryIdx,
	})
}

func (e *ModifyEvent) Resolve(es *EventSystem, effect EffectAction, context EffectContext, summaries *[]ActionSummary) EffectContext {
	if !IntegrityCheck(e.action, e.state) {
		slog.Warn("IntegrityCheck failed.")
		return nil
	}
	context, tmpSummary := e.action.Modify(e.state, context)
	appendSummary(summaries, e.action, SUMMARY_TYPE_MODIFY, tmpSummary)
	if isNil(context) {
		return nil
	}

	for _, branch := range e.branch {
		if !IntegrityCheck(branch.action, branch.state) {
			slog.Warn("IntegrityCheck failed.")
			continue
		}
		branch.Resolve(es, effect, context)
	}

	if e.modifier != nil && IntegrityCheck(e.modifier.action, e.modifier.state) {
		context = e.modifier.Resolve(es, effect, context, summaries)
	}

	return context
}

func appendSummary(summaries *[]ActionSummary, action Action, st SummaryType, summary Summary) {
	*summaries = append(*summaries, ActionSummary{Action: action, Type: st, Data: summary})
}

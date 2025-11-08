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
	active        Event
	ResultEffects []struct {
		Action  action.Action
		Context action.EffectContext
	}
}

type Event interface {
	Parent() Event
	State() interface{}
	Action() action.Action
	base() *baseEvent
}

type baseEvent struct {
	parent   Event
	branch   []*EffectEvent
	modifier *ModifyEvent
	state    *interface{}
}

type EffectEvent struct {
	baseEvent
	action action.EffectAction
}

type ModifyEvent struct {
	baseEvent
	action action.ModifyAction
}

func (e *baseEvent) Parent() Event {
	return e.parent
}

func (e *baseEvent) State() interface{} {
	return *e.state
}

func (e *EffectEvent) Action() action.Action {
	return e.action
}

func (e *ModifyEvent) Action() action.Action {
	return e.action
}

func (e *EffectEvent) base() *baseEvent {
	return &e.baseEvent
}

func (e *ModifyEvent) base() *baseEvent {
	return &e.baseEvent
}

func (es *EventSystem) listen(maybe interface{}, event Event) {
	if listener, ok := maybe.(Listener); ok {
		listener.Listen(es, event.Action(), event.State())
	}
}

func NewEvent(a action.Action, state *interface{}) Event {
	if effect, ok := a.(action.EffectAction); ok {
		return &EffectEvent{
			baseEvent: baseEvent{
				state: state,
			},
			action: effect,
		}
	} else if modifier, ok := a.(action.ModifyAction); ok {
		return &ModifyEvent{
			baseEvent: baseEvent{
				state: state,
			},
			action: modifier,
		}
	}
	return nil
}

func (es *EventSystem) Emit(event Event) {
	es.active = event
	players := es.Board.Players()
	for i, _ := range players {
		es.listen(players[(i+es.Board.Turn())%len(players)], event)
	}

	var searchStart, searchEnd, delta pkg.Point
	if es.Board.Turn()%2 == 0 {
		searchStart = pkg.Point{0, 0}
		searchEnd = pkg.SizeToPoint(es.Board.Size())
		delta = pkg.Point{1, 1}
	} else {
		searchStart = pkg.SizeToPoint(es.Board.Size())
		searchEnd = pkg.Point{0, 0}
		delta = pkg.Point{-1, -1}
	}

	for x := searchStart.X; x != searchEnd.X+delta.X; x += delta.X {
		for y := searchStart.Y; y != searchEnd.Y+delta.Y; y += delta.Y {
			es.listen(es.Board.Entities()[x][y], event)
		}
	}

	for x := searchStart.X; x != searchEnd.X+delta.X; x += delta.X {
		for y := searchStart.Y; y != searchEnd.Y+delta.Y; y += delta.Y {
			es.listen(es.Board.Structures()[x][y], event)
		}
	}

	if multi, ok := event.Action().(action.MultiEffectAction); ok {
		effects := multi.SubEffects(event.State())
		events := make([]Event, len(effects))
		for i, _ := range effects {
			events[i] = &effects[i]
		}
		es.ChainLine(events[0], events[1:])
	}
}

func (es *EventSystem) Chain(event Event) {
	if es.active == nil {
		slog.Error("prevActive is nil")
		return
	}
	prevActive := es.active
	defer func() { es.active = prevActive }()
	base := event.base()
	base.parent = prevActive
	if effect, ok := event.(*EffectEvent); ok {
		prevActive.base().branch = append(prevActive.base().branch, effect)
	} else if modifier, ok := event.(*ModifyEvent); ok {
		m := prevActive.base()
		for m.modifier != nil {
			m = m.modifier.base()
		}
		m.modifier = modifier
	}
	es.Emit(event)
}

func (es *EventSystem) ChainLine(event Event, pending []Event) {
	prevActive := es.active
	defer func() { es.active = prevActive }()
	base := event.base()
	base.parent = prevActive
	if effect, ok := event.(*EffectEvent); ok {
		prevActive.base().branch = append(prevActive.base().branch, effect)
	} else if modifier, ok := event.(*ModifyEvent); ok {
		m := prevActive.base()
		for m.modifier != nil {
			m = m.modifier.base()
		}
		m.modifier = modifier
	}
	es.Emit(event)
	if pending != nil && len(pending) > 0 {
		es.ChainLine(pending[0], pending[1:])
	}
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

func (e *EffectEvent) Resolve(es *EventSystem, beforeEffect action.EffectAction, beforeContext action.EffectContext) {
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
			Action  action.Action
			Context action.EffectContext
		}{Action: e.action, Context: context})
	}
}

func (e *ModifyEvent) Resolve(es *EventSystem, effect action.EffectAction, context action.EffectContext) action.EffectContext {
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

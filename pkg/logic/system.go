package logic

import (
	"frontage/pkg"
	"frontage/pkg/model"
	"log/slog"
	"reflect"
)

type Listener interface {
	Listen(es *EventSystem, event Action, state interface{})
}

type EventSystem struct {
	Board         *model.Board
	Trigger       *EffectEvent
	active        Event
	ResultEffects []struct {
		Action  Action
		Context EffectContext
	}
}

func (es *EventSystem) ListenHelper(maybe interface{}, action Action, state interface{}) {
	if listener, ok := maybe.(Listener); ok {
		listener.Listen(es, action, state)
	}
}

func (es *EventSystem) transmission(event Event) {
	es.active = event
	players := es.Board.Players()
	for i := range players {
		es.ListenHelper(players[(i+es.Board.Turn())%len(players)], event.Action(), event.State())
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
			es.ListenHelper(es.Board.Entities()[x][y], event.Action(), event.State())
		}
	}

	for x := searchStart.X; x != searchEnd.X+delta.X; x += delta.X {
		for y := searchStart.Y; y != searchEnd.Y+delta.Y; y += delta.Y {
			es.ListenHelper(es.Board.Structures()[x][y], event.Action(), event.State())
		}
	}

	if multi, ok := event.Action().(MultiEffectAction); ok {
		effects := multi.SubEffects(event.State())
		events := make([]Event, len(effects))
		for i := range effects {
			events[i] = effects[i]
		}
		es.ChainLine(events[0], events[1:]...)
	}
}

func (es *EventSystem) Emit(event *EffectEvent) {
	es.Trigger = event
	es.transmission(event)
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
	es.transmission(event)
}

func (es *EventSystem) ChainLine(event Event, pending ...Event) {
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
	es.transmission(event)
	if pending != nil && len(pending) > 0 {
		es.ChainLine(pending[0], pending[1:]...)
	}
}

func IntegrityCheck(a Action, state interface{}) bool {
	if a == nil {
		return false
	}
	if reflect.TypeOf(state) != a.WantState() {
		return false
	}
	return true
}

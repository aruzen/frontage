package logic

import (
	"frontage/pkg"
	"frontage/pkg/engine/model"
	"log/slog"
	"reflect"
)

type PreListener interface {
	PreListen(es *EventSystem, event Action, state interface{})
}

type RunningListener interface {
	RunningListen(es *EventSystem, event Action, state interface{})
}

type EventSystem struct {
	Board          *model.Board
	trigger        *EffectEvent
	active         Event
	isRunning      bool
	AppliedEffects []ActionResult
	Summaries      []*[]ActionSummary
}

func PreListenHelper(es *EventSystem, maybe interface{}, action Action, state interface{}) {
	if listener, ok := maybe.(PreListener); ok {
		listener.PreListen(es, action, state)
	}
}

func RunningListenHelper(es *EventSystem, maybe interface{}, action Action, state interface{}) {
	if listener, ok := maybe.(RunningListener); ok {
		listener.RunningListen(es, action, state)
	}
}

func (es *EventSystem) transmission(listenHelper func(es *EventSystem, maybe interface{}, action Action, state interface{}), event Event) {
	es.active = event
	players := es.Board.Players()
	for i := range players {
		listenHelper(es, players[(i+es.Board.Turn())%len(players)], event.Action(), event.State())
	}

	var searchStart, searchEnd, delta pkg.Point
	if es.Board.Turn()%2 == 0 {
		searchStart = pkg.Point{0, 0}
		searchEnd = pkg.SizeToPoint(es.Board.Size())
		delta = pkg.Point{1, 1}
	} else {
		searchStart = pkg.Point{es.Board.Size().Width - 1, es.Board.Size().Height - 1}
		searchEnd = pkg.Point{-1, -1}
		delta = pkg.Point{-1, -1}
	}

	for y := searchStart.Y; y != searchEnd.Y; y += delta.Y {
		for x := searchStart.X; x != searchEnd.X; x += delta.X {
			listenHelper(es, es.Board.Entities()[x][y], event.Action(), event.State())
		}
	}

	for y := searchStart.Y; y != searchEnd.Y; y += delta.Y {
		for x := searchStart.X; x != searchEnd.X; x += delta.X {
			listenHelper(es, es.Board.Structures()[x][y], event.Action(), event.State())
		}
	}

	if multi, ok := event.Action().(MultiEffectAction); ok {
		effects := multi.SubEffects(event.State())
		if len(effects) == 0 {
			return
		}
		events := make([]Event, len(effects))
		for i := range effects {
			events[i] = effects[i]
		}
		es.ChainLine(events[0], events[1:]...)
	}
}

func (es *EventSystem) Emit(event *EffectEvent) {
	es.trigger = event
	es.transmission(PreListenHelper, event)
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
	if !es.isRunning {
		es.transmission(PreListenHelper, event)
	} else {
		es.transmission(RunningListenHelper, event)
	}
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
	if !es.isRunning {
		es.transmission(PreListenHelper, event)
	} else {
		es.transmission(RunningListenHelper, event)
	}
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

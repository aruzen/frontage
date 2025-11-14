package logic

type Event interface {
	Parent() Event
	State() interface{}
	Action() Action
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
	action EffectAction
}

type ModifyEvent struct {
	baseEvent
	action ModifyAction
}

func (e *baseEvent) Parent() Event {
	return e.parent
}

func (e *baseEvent) State() interface{} {
	return *e.state
}

func (e *EffectEvent) Action() Action {
	return e.action
}

func (e *ModifyEvent) Action() Action {
	return e.action
}

func (e *EffectEvent) base() *baseEvent {
	return &e.baseEvent
}

func (e *ModifyEvent) base() *baseEvent {
	return &e.baseEvent
}

func NewEvent(a Action, state *interface{}) Event {
	if effect, ok := a.(EffectAction); ok {
		return &EffectEvent{
			baseEvent: baseEvent{
				state: state,
			},
			action: effect,
		}
	} else if modifier, ok := a.(ModifyAction); ok {
		return &ModifyEvent{
			baseEvent: baseEvent{
				state: state,
			},
			action: modifier,
		}
	}
	return nil
}

func NewEffectEvent(a EffectAction, state *interface{}) *EffectEvent {
	return &EffectEvent{
		baseEvent: baseEvent{
			state: state,
		},
		action: a,
	}
}

func NewModifyEvent(a ModifyAction, state *interface{}) *ModifyEvent {
	return &ModifyEvent{
		baseEvent: baseEvent{
			state: state,
		},
		action: a,
	}
}

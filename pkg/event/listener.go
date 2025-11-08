package event

import (
	"frontage/pkg/event/action"
)

type Listener interface {
	Listen(es *EventSystem, event action.Action, state interface{})
}

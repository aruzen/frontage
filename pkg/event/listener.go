package event

import (
	"frontage/pkg/event/action"
)

type Listener interface {
	Listen(event action.Action, state interface{})
}

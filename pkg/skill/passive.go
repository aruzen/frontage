package skill

import "frontage/pkg/event"

type PassiveSkill interface {
	event.Listener
}

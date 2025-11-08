package skill

import (
	"frontage/pkg"
	"frontage/pkg/event"
)

type Skill interface {
}

type NamedSkill interface {
	pkg.Localized
	Skill
}

type ActiveSkill interface {
	Skill
	Active(board *pkg.Board)
}

type PassiveSkill interface {
	Skill
	event.Listener
}

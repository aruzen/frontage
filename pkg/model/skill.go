package model

import (
	"frontage/pkg"
)

type Skill interface {
}

type NamedSkill interface {
	pkg.Localized
	Skill
}

type ActiveSkill interface {
	Skill
	Active(board *Board)
}

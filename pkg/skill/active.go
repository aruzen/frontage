package skill

import "frontage/pkg"

type ActiveSkill interface {
	Active(board *pkg.Board)
}

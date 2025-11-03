package entity

import "frontage/pkg"

type Entity interface {
	HP() int
	MP() int
	ATK() int
	LegalMoves() []pkg.Point

	SetHP(int)
	SetMP(int)
	SetATK(int)

	PassiveSkill()
	ActiveSkill()
}

package entity

import "frontage/pkg"

type Entity interface {
	HP() int
	MP() int
	ATK() int
	LegalMoves() []pkg.Point
	PassiveSkill()
	ActiveSkill()

	Copy() MutableEntity
}

type MutableEntity interface {
	Entity
	SetHP(int)
	SetMP(int)
	SetATK(int)
}

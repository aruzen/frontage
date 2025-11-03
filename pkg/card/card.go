package card

import "frontage/pkg"

type Placed int

const (
	PLACED_MAIN_DECK Placed = iota
	PLACED_SUB_DECK
	PLACED_PLAYER_HANDS
	PLACED_GRAVEYARD
	PLACED_EXTRA
)

type Type int

const (
	TYPE_PIECE Type = iota
	TYPE_CHARM
	TYPE_STRUCTURE
	TYPE_MAGIC
	TYPE_SKILL
)

type Card interface {
	Name() string
	Resource() string
	Placed() Placed
	Type() Type
	PlayCost() pkg.Materials

	CardCopy() MutableCard
}

type MutableCard interface {
	Card

	SetPlaced(Placed)
	SetPlayCost(playCost pkg.Materials)
}

type BaseCard struct {
	name     string
	resource string
	placed   Placed
	playCost pkg.Materials
}

func (b *BaseCard) CardCopy() MutableCard {
	return &BaseCard{b.name, b.resource, b.placed, b.playCost.Copy()}
}

func (b *BaseCard) Name() string {
	return b.name
}

func (b *BaseCard) Resource() string {
	return b.resource
}

func (b *BaseCard) Placed() Placed {
	return b.placed
}

func (b *BaseCard) PlayCost() pkg.Materials {
	return b.playCost.Copy()
}

func (b *BaseCard) SetPlaced(placed Placed) {
	b.placed = placed
}

func (b *BaseCard) SetPlayCost(playCost pkg.Materials) {
	b.playCost = playCost
}

func (b *BaseCard) Type() Type {
	panic("オーバーライドしてください")
}

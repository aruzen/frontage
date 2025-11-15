package model

import (
	"github.com/google/uuid"
)

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
	Id() uuid.UUID
	Name() string
	Resource() string
	Placed() Placed
	Type() Type
	PlayCost() Materials

	CardCopy() MutableCard
	Mirror(i uuid.UUID) MutableCard
}

type MutableCard interface {
	Card

	SetPlaced(Placed)
	SetPlayCost(playCost Materials)
}

type BaseCard struct {
	id       uuid.UUID
	name     string
	resource string
	placed   Placed
	playCost Materials
}

var _ MutableCard = (*BaseCard)(nil)

func NewBaseCard(name, resource string, placed Placed, playCost Materials) *BaseCard {
	return &BaseCard{
		id:       uuid.New(),
		name:     name,
		resource: resource,
		placed:   placed,
		playCost: playCost.Copy(),
	}
}

func (b *BaseCard) CardCopy() MutableCard {
	return &BaseCard{
		id:       b.id,
		name:     b.name,
		resource: b.resource,
		placed:   b.placed,
		playCost: b.playCost.Copy(),
	}
}

func (b *BaseCard) Mirror(i uuid.UUID) MutableCard {
	copied := b.CardCopy().(*BaseCard)
	copied.id = i
	return copied
}

func (b *BaseCard) Id() uuid.UUID {
	return b.id
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

func (b *BaseCard) PlayCost() Materials {
	return b.playCost.Copy()
}

func (b *BaseCard) SetPlaced(placed Placed) {
	b.placed = placed
}

func (b *BaseCard) SetPlayCost(playCost Materials) {
	b.playCost = playCost
}

func (b *BaseCard) Type() Type {
	panic("オーバーライドしてください")
}

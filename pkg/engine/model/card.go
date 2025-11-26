package model

import (
	"frontage/pkg"
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

type CardType int

const (
	CARD_TYPE_PIECE CardType = iota
	CARD_TYPE_CHARM
	CARD_TYPE_STRUCTURE
	CARD_TYPE_MAGIC
	CARD_TYPE_SKILL
)

type Card interface {
	pkg.Localized
	Id() uuid.UUID
	Placed() Placed
	Type() CardType
	PlayCost() Materials

	CardCopy() MutableCard
	CardMirror(i uuid.UUID) MutableCard
}

type MutableCard interface {
	Card

	SetPlaced(Placed)
	SetPlayCost(playCost Materials)
}

type BaseCard struct {
	id       uuid.UUID
	tag      pkg.LocalizeTag
	placed   Placed
	playCost Materials
}

var _ MutableCard = (*BaseCard)(nil)

func NewBaseCard(tag pkg.LocalizeTag, placed Placed, playCost Materials) *BaseCard {
	return &BaseCard{
		id:       uuid.New(),
		tag:      tag,
		placed:   placed,
		playCost: playCost.Copy(),
	}
}

func (b *BaseCard) CardCopy() MutableCard {
	return &BaseCard{
		id:       b.id,
		tag:      b.tag,
		placed:   b.placed,
		playCost: b.playCost.Copy(),
	}
}

func (b *BaseCard) CardMirror(i uuid.UUID) MutableCard {
	copied := b.CardCopy().(*BaseCard)
	copied.id = i
	return copied
}

func (b *BaseCard) Id() uuid.UUID {
	return b.id
}

func (b *BaseCard) LocalizeTag() pkg.LocalizeTag {
	return b.tag
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

func (b *BaseCard) Type() CardType {
	panic("オーバーライドしてください")
}

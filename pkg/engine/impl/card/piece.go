package card

import (
	"frontage/pkg"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

type Piece interface {
	model.Card
	HP() int
	MP() int
	ATK() int
	// LegalMoves 返り値を変更しないこと
	LegalMoves() []pkg.Point
	// AttackRanges 返り値を変更しないこと
	AttackRanges() []pkg.Point

	Copy() MutablePiece
	Mirror(uuid uuid.UUID) MutablePiece
}

type MutablePiece interface {
	Piece
	model.MutableCard
	SetHP(int)
	SetMP(int)
	SetATK(int)
}

type BasePiece struct {
	model.BaseCard
	legalMoves   []pkg.Point
	attackRanges []pkg.Point
	hp, mp, atk  int
}

func NewBasePiece(base *model.BaseCard, hp, mp, atk int, legalMoves, attackRanges []pkg.Point) *BasePiece {
	return &BasePiece{
		BaseCard:     *base,
		legalMoves:   legalMoves,
		attackRanges: attackRanges,
		hp:           hp,
		mp:           mp,
		atk:          atk,
	}
}

func (p *BasePiece) Type() model.CardType {
	return model.CARD_TYPE_PIECE
}

func (p *BasePiece) HP() int {
	return p.hp
}

func (p *BasePiece) MP() int {
	return p.mp
}

func (p *BasePiece) ATK() int {
	return p.atk
}

func (p *BasePiece) LegalMoves() []pkg.Point {
	result := make([]pkg.Point, len(p.legalMoves))
	copy(result, p.legalMoves)
	return result
}

func (p *BasePiece) AttackRanges() []pkg.Point {
	result := make([]pkg.Point, len(p.attackRanges))
	copy(result, p.attackRanges)
	return result
}

func (p *BasePiece) Copy() MutablePiece {
	legalMoves := make([]pkg.Point, len(p.legalMoves))
	copy(legalMoves, p.legalMoves)
	attackRanges := make([]pkg.Point, len(p.attackRanges))
	copy(attackRanges, p.attackRanges)
	return &BasePiece{
		BaseCard:     *p.BaseCard.CardCopy().(*model.BaseCard),
		legalMoves:   legalMoves,
		attackRanges: attackRanges,
		hp:           p.hp,
		mp:           p.mp,
		atk:          p.atk,
	}
}

func (p *BasePiece) Mirror(uuid uuid.UUID) MutablePiece {
	legalMoves := make([]pkg.Point, len(p.legalMoves))
	copy(legalMoves, p.legalMoves)
	attackRanges := make([]pkg.Point, len(p.attackRanges))
	copy(attackRanges, p.attackRanges)
	return &BasePiece{
		BaseCard:     *p.BaseCard.CardMirror(uuid).(*model.BaseCard),
		legalMoves:   legalMoves,
		attackRanges: attackRanges,
		hp:           p.hp,
		mp:           p.mp,
		atk:          p.atk,
	}
}

func (p *BasePiece) CardCopy() model.MutableCard {
	return p.Copy().(model.MutableCard)
}

func (p *BasePiece) SetHP(i int) {
	p.hp = i
}

func (p *BasePiece) SetMP(i int) {
	p.mp = i
}

func (p *BasePiece) SetATK(i int) {
	p.atk = i
}

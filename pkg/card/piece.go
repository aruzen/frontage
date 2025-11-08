package card

import "frontage/pkg"

type Piece interface {
	Card
	HP() int
	MP() int
	ATK() int
	// LegalMoves 返り値を変更しないこと
	LegalMoves() []pkg.Point
	// AttackRanges 返り値を変更しないこと
	AttackRanges() []pkg.Point

	Copy() MutablePiece
}

type MutablePiece interface {
	Piece
	SetHP(int)
	SetMP(int)
	SetATK(int)
}

type BasePiece struct {
	BaseCard
	legalMoves   []pkg.Point
	attackRanges []pkg.Point
	hp, mp, atk  int
}

func (p *BasePiece) Type() Type {
	return TYPE_PIECE
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
	return p.legalMoves
}

func (p *BasePiece) AttackRanges() []pkg.Point {
	return p.attackRanges
}

func (p *BasePiece) Copy() MutablePiece {
	return &BasePiece{
		BaseCard: BaseCard{
			id:       p.id,
			name:     p.name,
			resource: p.resource,
			placed:   p.placed,
			playCost: p.playCost.Copy(),
		},
		legalMoves:   p.legalMoves,
		attackRanges: p.attackRanges,
		hp:           p.hp,
		mp:           p.mp,
		atk:          p.atk,
	}
}

func (p *BasePiece) CardCopy() MutableCard {
	return p.Copy().(MutableCard)
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

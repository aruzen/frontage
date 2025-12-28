package model

import (
	"frontage/pkg"
	"github.com/google/uuid"
)

type Piece interface {
	Id() uuid.UUID
	Position() pkg.Point
	HP() int
	MP() int
	ATK() int
	Owner() Player

	HaveSkill(s SkillTag) bool
	Skills() []Skill
	LegalMoves() []pkg.Point
	AttackRanges() []pkg.Point

	Copy() MutablePiece
	Mirror() MutablePiece
}

type MutablePiece interface {
	Piece

	GetSkill(s SkillTag) (Skill, bool)
	AddSkill(s Skill)
	RemoveSkill(s Skill)
	SetHP(int)
	SetMP(int)
	SetATK(int)

	// SetPosition この関数ではボード側にも適用されないので、呼び出さないでください。`Board.SetPiece(Piece)`を使用してください。
	SetPosition(pos pkg.Point)
}

type BasePiece struct {
	id       uuid.UUID
	position pkg.Point
	hp       int
	mp       int
	atk      int
	owner    Player

	skills       []Skill
	legalMoves   []pkg.Point
	attackRanges []pkg.Point
	actedCount   int
	maxActions   int
}

var _ MutablePiece = (*BasePiece)(nil)

func NewBasePiece(id uuid.UUID, owner Player, pos pkg.Point, hp, mp, atk int, legalMoves, attackRanges []pkg.Point, skills ...Skill) *BasePiece {
	piece := &BasePiece{
		id:           id,
		owner:        owner,
		position:     pos,
		hp:           hp,
		mp:           mp,
		atk:          atk,
		legalMoves:   append([]pkg.Point(nil), legalMoves...),
		attackRanges: append([]pkg.Point(nil), attackRanges...),
		maxActions:   1,
	}
	if len(skills) > 0 {
		piece.skills = append([]Skill(nil), skills...)
	}
	return piece
}

func (e *BasePiece) Position() pkg.Point {
	return e.position
}

func (e *BasePiece) Id() uuid.UUID {
	return e.id
}

func (e *BasePiece) HP() int {
	return e.hp
}

func (e *BasePiece) MP() int {
	return e.mp
}

func (e *BasePiece) ATK() int {
	return e.atk
}

func (e *BasePiece) Owner() Player {
	return e.owner
}

func (e *BasePiece) HaveSkill(s SkillTag) bool {
	_, ok := e.GetSkill(s)
	return ok
}

func (e *BasePiece) Skills() []Skill {
	result := make([]Skill, len(e.skills))
	copy(result, e.skills)
	return result
}

// LegalMoves は返り値を書き換えないでください。
func (e *BasePiece) LegalMoves() []pkg.Point {
	result := make([]pkg.Point, len(e.legalMoves))
	copy(result, e.legalMoves)
	return result
}

// AttackRanges は返り値を書き換えないでください。
func (e *BasePiece) AttackRanges() []pkg.Point {
	result := make([]pkg.Point, len(e.attackRanges))
	copy(result, e.attackRanges)
	return result
}

func (e *BasePiece) Copy() MutablePiece {
	copySkills := make([]Skill, len(e.skills))
	copy(copySkills, e.skills)
	return &BasePiece{
		id:           e.id,
		position:     e.position,
		hp:           e.hp,
		mp:           e.mp,
		atk:          e.atk,
		owner:        e.owner,
		legalMoves:   append([]pkg.Point(nil), e.legalMoves...),
		attackRanges: append([]pkg.Point(nil), e.attackRanges...),
		skills:       copySkills,
		actedCount:   e.actedCount,
		maxActions:   e.maxActions,
	}
}

func (e *BasePiece) Mirror() MutablePiece {
	copied := e.Copy().(*BasePiece)
	copied.id = uuid.New()
	return copied
}

func (e *BasePiece) GetSkill(tag SkillTag) (Skill, bool) {
	for _, sk := range e.skills {
		if sk.Tag() == tag {
			return sk, true
		}
	}
	return nil, false
}

func (e *BasePiece) AddSkill(s Skill) {
	if s == nil {
		return
	}
	e.skills = append(e.skills, s)
}

func (e *BasePiece) RemoveSkill(target Skill) {
	for i, sk := range e.skills {
		if sk == target {
			e.skills = append(e.skills[:i], e.skills[i+1:]...)
			return
		}
	}
}

func (e *BasePiece) SetHP(v int) {
	e.hp = v
}

func (e *BasePiece) SetMP(v int) {
	e.mp = v
}

func (e *BasePiece) SetATK(v int) {
	e.atk = v
}

func (e *BasePiece) SetPosition(pos pkg.Point) {
	e.position = pos
}

func (e *BasePiece) ActionsUsedThisTurn() int {
	return e.actedCount
}

func (e *BasePiece) MaxActionsPerTurn() int {
	return e.maxActions
}

func (e *BasePiece) SetActionsUsedThisTurn(v int) {
	e.actedCount = v
	if e.actedCount < 0 {
		e.actedCount = 0
	}
}

func (e *BasePiece) SetMaxActionsPerTurn(v int) {
	if v < 0 {
		v = 0
	}
	e.maxActions = v
	if e.actedCount > e.maxActions {
		e.actedCount = e.maxActions
	}
}

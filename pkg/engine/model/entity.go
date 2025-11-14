package model

import (
	"frontage/pkg/engine"
	"github.com/google/uuid"
)

type Entity interface {
	Id() uuid.UUID
	Position() pkg.Point
	HP() int
	MP() int
	ATK() int
	Owner() *Player

	HaveSkill(s NamedSkill) bool
	Skills() []Skill
	LegalMoves() []pkg.Point
	AttackRanges() []pkg.Point

	Copy() MutableEntity
	Mirror() MutableEntity
}

type MutableEntity interface {
	Entity

	GetSkill(s pkg.LocalizeTag) (NamedSkill, bool)
	AddSkill(s Skill)
	RemoveSkill(s Skill)
	SetHP(int)
	SetMP(int)
	SetATK(int)

	// SetPosition この関数ではボード側にも適用されないので、呼び出さないでください。`Board.SetEntity(Entity)`を使用してください。
	SetPosition(pos pkg.Point)
}

type BaseEntity struct {
	id       uuid.UUID
	position pkg.Point
	hp       int
	mp       int
	atk      int
	owner    *Player

	skills       []Skill
	legalMoves   []pkg.Point
	attackRanges []pkg.Point
}

var _ MutableEntity = (*BaseEntity)(nil)

func NewBaseEntity(owner *Player, pos pkg.Point, hp, mp, atk int, legalMoves, attackRanges []pkg.Point, skills ...Skill) *BaseEntity {
	entity := &BaseEntity{
		id:           uuid.New(),
		owner:        owner,
		position:     pos,
		hp:           hp,
		mp:           mp,
		atk:          atk,
		legalMoves:   append([]pkg.Point(nil), legalMoves...),
		attackRanges: append([]pkg.Point(nil), attackRanges...),
	}
	if len(skills) > 0 {
		entity.skills = append([]Skill(nil), skills...)
	}
	return entity
}

func (e *BaseEntity) Position() pkg.Point {
	return e.position
}

func (e *BaseEntity) Id() uuid.UUID {
	return e.id
}

func (e *BaseEntity) HP() int {
	return e.hp
}

func (e *BaseEntity) MP() int {
	return e.mp
}

func (e *BaseEntity) ATK() int {
	return e.atk
}

func (e *BaseEntity) Owner() *Player {
	return e.owner
}

func (e *BaseEntity) HaveSkill(s NamedSkill) bool {
	if s == nil {
		return false
	}
	_, ok := e.GetSkill(s.Name())
	return ok
}

func (e *BaseEntity) Skills() []Skill {
	result := make([]Skill, len(e.skills))
	copy(result, e.skills)
	return result
}

// LegalMoves は返り値を書き換えないでください。
func (e *BaseEntity) LegalMoves() []pkg.Point {
	result := make([]pkg.Point, len(e.legalMoves))
	copy(result, e.legalMoves)
	return result
}

// AttackRanges は返り値を書き換えないでください。
func (e *BaseEntity) AttackRanges() []pkg.Point {
	result := make([]pkg.Point, len(e.attackRanges))
	copy(result, e.attackRanges)
	return result
}

func (e *BaseEntity) Copy() MutableEntity {
	copySkills := make([]Skill, len(e.skills))
	copy(copySkills, e.skills)
	return &BaseEntity{
		id:           e.id,
		position:     e.position,
		hp:           e.hp,
		mp:           e.mp,
		atk:          e.atk,
		owner:        e.owner,
		legalMoves:   append([]pkg.Point(nil), e.legalMoves...),
		attackRanges: append([]pkg.Point(nil), e.attackRanges...),
		skills:       copySkills,
	}
}

func (e *BaseEntity) Mirror() MutableEntity {
	copied := e.Copy().(*BaseEntity)
	copied.id = uuid.New()
	return copied
}

func (e *BaseEntity) GetSkill(tag pkg.LocalizeTag) (NamedSkill, bool) {
	for _, sk := range e.skills {
		if named, ok := sk.(NamedSkill); ok && named.Name() == tag {
			return named, true
		}
	}
	return nil, false
}

func (e *BaseEntity) AddSkill(s Skill) {
	if s == nil {
		return
	}
	e.skills = append(e.skills, s)
}

func (e *BaseEntity) RemoveSkill(target Skill) {
	for i, sk := range e.skills {
		if sk == target {
			e.skills = append(e.skills[:i], e.skills[i+1:]...)
			return
		}
	}
}

func (e *BaseEntity) SetHP(v int) {
	e.hp = v
}

func (e *BaseEntity) SetMP(v int) {
	e.mp = v
}

func (e *BaseEntity) SetATK(v int) {
	e.atk = v
}

func (e *BaseEntity) SetPosition(pos pkg.Point) {
	e.position = pos
}

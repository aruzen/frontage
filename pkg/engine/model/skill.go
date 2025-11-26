package model

type SkillTag string

type Skill interface {
	Tag() SkillTag
}

type ActiveSkill interface {
	Skill
	Active(board *Board)
}

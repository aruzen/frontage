package data

type Skill struct {
	Name      string
	Data      map[string]interface{}
	Subskills []Skill
}

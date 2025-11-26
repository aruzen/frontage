package data

// TODO: implement me

type Skill struct {
	Tag       string                 `json:"tag"`
	Data      map[string]interface{} `json:"data"`
	Subskills []Skill                `json:"subskills"`
}

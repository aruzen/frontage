package data

import (
	"github.com/google/uuid"
)

type Card struct {
	Placed int       `json:"placed"`
	Type   int       `json:"type"`
	Tag    string    `json:"tag"`
	UUID   uuid.UUID `json:"uuid"`
	Cost   Materials `json:"cost"`
	Skills []Skill   `json:"skills"`
}

type PieceCard struct {
	Card
	Atk int `json:"atk"`
	Hp  int `json:"hp"`
	Mp  int `json:"mp"`
}

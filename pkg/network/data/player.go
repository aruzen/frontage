package data

import "github.com/google/uuid"

type Player struct {
	UUID     uuid.UUID
	Name     string
	MainDeck []Card
	SubDeck  []Card
}

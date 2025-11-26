package data

import "github.com/google/uuid"

type Player struct {
	Id       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	MainDeck []string  `json:"main_deck"`
	SubDeck  []string  `json:"sub_deck"`
}

package game_api

import "github.com/google/uuid"

type GameInitializePacket struct {
	Width    int `json:"width"`
	Height   int `json:"height"`
	YourSide int `json:"your_side"`
}

type OpponentPlayerInitializePacket struct {
	Id   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

type MyDeckUploadPacket struct {
	Id       uuid.UUID `json:"id"`
	MainDeck []string  `json:"main_deck"`
	SubDeck  []string  `json:"sub_deck"`
}

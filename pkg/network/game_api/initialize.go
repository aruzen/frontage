package game_api

import (
	"frontage/pkg/network"
	"github.com/google/uuid"
)

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

func (p GameInitializePacket) PacketTag() network.PacketTag {
	return network.GAME_INITIALIZE_PACKET_TAG
}

func (p OpponentPlayerInitializePacket) PacketTag() network.PacketTag {
	return network.OPPONENT_PLAYER_INITIALIZE_PACKET_TAG
}

func (p MyDeckUploadPacket) PacketTag() network.PacketTag {
	return network.MY_DECK_UPLOAD_PACKET_TAG
}

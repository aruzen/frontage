package game_dispatcher

import (
	"encoding/json"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/game_api"
	"github.com/google/uuid"
)

type GameInitializeDispatcher struct {
}

func NewGameInitializePacket() *GameInitializeDispatcher {
	return &GameInitializeDispatcher{}
}

func (d *GameInitializeDispatcher) DispatchPacket(b *model.Board, p model.Player) ([]byte, error) {
	packet := game_api.GameInitializePacket{
		Width:  b.Size().Width,
		Height: b.Size().Height,
	}

	if players := b.Players(); players[0] != nil && players[0].ID() == p.ID() {
		packet.YourSide = 0
	} else {
		packet.YourSide = 1
	}
	return json.Marshal(packet)
}

type OpponentPlayerInitializeDispatcher struct{}

func NewOpponentPlayerInitializeDispatcher() *OpponentPlayerInitializeDispatcher {
	return &OpponentPlayerInitializeDispatcher{}
}

func (d *OpponentPlayerInitializeDispatcher) DispatchPacket(id uuid.UUID, name string) ([]byte, error) {
	packet := game_api.OpponentPlayerInitializePacket{
		Id:   id,
		Name: name,
	}
	return json.Marshal(packet)
}

type MyDeckUploadDispatcher struct{}

func NewMyDeckUploadDispatcher() *MyDeckUploadDispatcher {
	return &MyDeckUploadDispatcher{}
}

func (d *MyDeckUploadDispatcher) DispatchPacket(id uuid.UUID, mainDeck, subDeck []string) ([]byte, error) {
	packet := game_api.MyDeckUploadPacket{
		Id:       id,
		MainDeck: mainDeck,
		SubDeck:  subDeck,
	}
	return json.Marshal(packet)
}

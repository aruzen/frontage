package game_dispatcher

import (
	"frontage/pkg/network"
	"frontage/pkg/network/game_api"
)

type TurnStartDispatcher struct{}

type TurnEndDispatcher struct{}

func NewTurnStartDispatcher() *TurnStartDispatcher {
	return &TurnStartDispatcher{}
}

func NewTurnEndDispatcher() *TurnEndDispatcher {
	return &TurnEndDispatcher{}
}

func (d *TurnStartDispatcher) DispatchPacket(turn int) (network.Packet, error) {
	return game_api.TurnStartPacket{Turn: turn}, nil
}

func (d *TurnEndDispatcher) DispatchPacket(turn int) (network.Packet, error) {
	return game_api.TurnEndPacket{Turn: turn}, nil
}

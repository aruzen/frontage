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

func (d *TurnStartDispatcher) DispatchPacket(packet game_api.TurnStartPacket) (network.Packet, error) {
	return packet, nil
}

func (d *TurnEndDispatcher) DispatchPacket(packet game_api.TurnEndPacket) (network.Packet, error) {
	return packet, nil
}

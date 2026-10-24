package game_handler

import (
	"encoding/json"
	"frontage/pkg/network/game_api"
)

type TurnPassHandler struct{}
type TurnStartHandler struct{}
type TurnEndHandler struct{}

func NewTurnPassHandler() *TurnPassHandler {
	return &TurnPassHandler{}
}

func (h *TurnPassHandler) ServePacket(data []byte) (int, error) {
	var packet game_api.TurnPassPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return 0, err
	}
	return packet.Turn, nil
}

func NewTurnStartHandler() *TurnStartHandler {
	return &TurnStartHandler{}
}

func (h *TurnStartHandler) ServePacket(data []byte) (game_api.TurnStartPacket, error) {
	var packet game_api.TurnStartPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return game_api.TurnStartPacket{}, err
	}
	return packet, nil
}

func NewTurnEndHandler() *TurnEndHandler {
	return &TurnEndHandler{}
}

func (h *TurnEndHandler) ServePacket(data []byte) (game_api.TurnEndPacket, error) {
	var packet game_api.TurnEndPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return game_api.TurnEndPacket{}, err
	}
	return packet, nil
}

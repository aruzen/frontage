package game_handler

import (
	"encoding/json"
	"frontage/pkg/network/game_api"
	"frontage/pkg/network/repository"
)

// ActEventHandler parses incoming action event packets.
type ActEventHandler struct {
	actionRepo *repository.ActionRepository
}

func NewActEventHandler(actionRepo *repository.ActionRepository) *ActEventHandler {
	return &ActEventHandler{
		actionRepo: actionRepo,
	}
}

func (h *ActEventHandler) ServePacket(data []byte) (game_api.ActEventPacket, error) {
	var packet game_api.ActEventPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return game_api.ActEventPacket{}, err
	}
	return packet, nil
}

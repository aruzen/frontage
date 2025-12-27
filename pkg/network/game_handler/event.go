package game_handler

import (
	"encoding/json"
	"frontage/pkg/network/game_api"
)

// ActEventHandler parses incoming action event packets.
type ActEventHandler struct{}

func NewActEventHandler() *ActEventHandler { return &ActEventHandler{} }

func (h *ActEventHandler) ServePacket(data []byte) ([]game_api.ActEventPayload, error) {
	var packet game_api.ActEventPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return nil, err
	}
	return packet.Events, nil
}

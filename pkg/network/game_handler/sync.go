package game_handler

import (
	"encoding/json"
	"frontage/pkg/network/game_api"
)

type SyncHandler struct{}

func NewSyncHandler() *SyncHandler {
	return &SyncHandler{}
}

func (h *SyncHandler) ServePacket(data []byte) (int, error) {
	var packet game_api.SyncPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return 0, err
	}
	return packet.SyncCause, nil
}

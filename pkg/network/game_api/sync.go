package game_api

import "frontage/pkg/network"

const (
	SYNC_CAUSE_GAME_START = iota
	SYNC_CAUSE_INITIALIZE_END
	SYNC_CAUSE_DRAW_PHASE_END
	SYNC_CAUSE_TURN_END
)

// SyncPacket is an empty packet for client/server sync.
type SyncPacket struct {
	SyncCause int `json:"sync_cause"`
}

func (p SyncPacket) PacketTag() network.PacketTag {
	return network.SYNC_PACKET_TAG
}

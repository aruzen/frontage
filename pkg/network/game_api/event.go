package game_api

import (
	"frontage/pkg/network"
	"frontage/pkg/network/data"
)

type ActEventPayload struct {
	Result  data.ActionResult
	Summary []data.ActionSummary
}

type ActEventPacket struct {
	Events []ActEventPayload
}

func (p ActEventPacket) PacketTag() network.PacketTag {
	return network.ACT_EVENT_PACKET_TAG
}

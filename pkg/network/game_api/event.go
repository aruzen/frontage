package game_api

import (
	"frontage/pkg/network"
	"frontage/pkg/network/data"
)

type ActEventPacket struct {
	Events    []data.ActionResult
	Summaries [][]data.ActionSummary
}

func (p ActEventPacket) PacketTag() network.PacketTag {
	return network.ACT_EVENT_PACKET_TAG
}

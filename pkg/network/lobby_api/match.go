package lobby_api

import (
	"frontage/pkg/network"
	"frontage/pkg/network/data"
	"github.com/google/uuid"
)

type WaitMatchMakePacket struct {
	Type data.MatchType
}

type CompleteMatchMakePacket struct {
	MatchID uuid.UUID
}

func (p WaitMatchMakePacket) PacketTag() network.PacketTag {
	return network.WAIT_MATCH_MAKE_PACKET_TAG
}

func (p CompleteMatchMakePacket) PacketTag() network.PacketTag {
	return network.COMPLETE_MATCH_MAKE_PACKET_TAG
}

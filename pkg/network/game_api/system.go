package game_api

import "frontage/pkg/network"

type TurnStartPacket struct {
	Turn       int  `json:"turn"`
	IsYourSide bool `json:"is_your_side"`
	TurnTime   int  `json:"turn_time"`
	ExtraTime  int  `json:"extra_time"`
}

type TurnEndPacket struct {
	Turn       int  `json:"turn"`
	IsYourSide bool `json:"is_your_side"`
}

type TurnPassPacket struct {
	Turn int `json:"turn"`
}

func (p TurnStartPacket) PacketTag() network.PacketTag {
	return network.TURN_START_PACKET_TAG
}

func (p TurnEndPacket) PacketTag() network.PacketTag {
	return network.TURN_END_PACKET_TAG
}

func (p TurnPassPacket) PacketTag() network.PacketTag {
	return network.TURN_PASS_PACKET_TAG
}

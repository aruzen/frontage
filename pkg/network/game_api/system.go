package game_api

import "frontage/pkg/network"

type TurnStartPacket struct {
	Turn int `json:"turn"`
}

type TurnEndPacket struct {
	Turn int `json:"turn"`
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

package game_api

import (
	"frontage/pkg"
	"frontage/pkg/network"
	"frontage/pkg/network/data"
)

type ActPieceSkillPacket struct {
	Position pkg.Point  `json:"position"`
	Skill    data.Skill `json:"skill"`
}

type ActStructureSkillPacket struct {
	Position pkg.Point  `json:"position"`
	Skill    data.Skill `json:"skill"`
}

type PlayCardPacket struct {
	Card data.Card `json:"card"`
}

func (p ActPieceSkillPacket) PacketTag() network.PacketTag {
	return network.ACT_PIECE_SKILL_PACKET_TAG
}

func (p ActStructureSkillPacket) PacketTag() network.PacketTag {
	return network.ACT_STRUCTURE_SKILL_PACKET_TAG
}

func (p PlayCardPacket) PacketTag() network.PacketTag {
	return network.PLAY_CARD_PACKET_TAG
}

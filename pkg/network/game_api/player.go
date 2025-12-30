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

const (
	DRAW_TYPE_MAIN_ONLY = iota
	DRAW_TYPE_SUB_ONLY
	DRAW_TYPE_MAIN_DRAW_SUB_OPEN
	DRAW_TYPE_MAIN_OPEN_SUB_DRAW
)

type DrawPhasePacket struct {
	DrawType int `json:"draw_type"`
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

func (p DrawPhasePacket) PacketTag() network.PacketTag {
	return network.DRAW_PHASE_PACKET_TAG
}

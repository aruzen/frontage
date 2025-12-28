package network

type PacketTag uint16

type Packet interface {
	PacketTag() PacketTag
}

type PacketHeader struct {
	PacketTag PacketTag
}

type UnsolvedPacket struct {
	Tag  PacketTag
	Body []byte
}

const (
	SystemPacketFlag PacketTag = iota
	LobbyPacketFlag
	GamePacketFlag
)

const (
	WAIT_MATCH_MAKE_PACKET_TAG PacketTag = iota + LobbyPacketFlag<<14
	COMPLETE_MATCH_MAKE_PACKET_TAG
)

const (
	ACT_EVENT_PACKET_TAG PacketTag = iota + GamePacketFlag<<14
	GAME_INITIALIZE_PACKET_TAG
	OPPONENT_PLAYER_INITIALIZE_PACKET_TAG
	MY_DECK_UPLOAD_PACKET_TAG
)

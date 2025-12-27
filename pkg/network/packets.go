package network

import (
	"encoding/binary"
	"encoding/json"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
)

type PacketTag uint16

type Packet interface {
	PacketTag() PacketTag
}

type PacketHeader struct {
	PacketTag PacketTag
}

func SendPacket(id uuid.UUID, h Packet) bool {
	connection := repository.GetConnection(id)
	if connection == nil {
		return false
	}
	body, err := json.Marshal(h)
	if err != nil {
		return false
	}
	header := make([]byte, 6)
	binary.LittleEndian.PutUint16(header[:2], uint16(h.PacketTag()))
	binary.LittleEndian.PutUint32(header[2:6], uint32(len(body)))
	connection.Mtx.Lock()
	defer connection.Mtx.Unlock()
	writeAll := func(data []byte) bool {
		for len(data) > 0 {
			written, err := connection.Conn.Write(data)
			if err != nil {
				return false
			}
			if written <= 0 {
				return false
			}
			data = data[written:]
		}
		return true
	}
	if !writeAll(header) {
		return false
	}
	if !writeAll(body) {
		return false
	}
	return true
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

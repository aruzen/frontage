package controller

import (
	"context"
	"frontage/pkg/network"
	"frontage/pkg/network/lobby_handler"
	"github.com/google/uuid"
)

type LobbyPacketHandlers struct {
	MatchMake *lobby_handler.MatchMakeHandler
}

func DispatchLobbyPacket(handlers LobbyPacketHandlers, tag network.PacketTag, id uuid.UUID, data []byte) error {
	return DispatchLobbyPacketContext(context.Background(), handlers, tag, id, data)
}

func DispatchLobbyPacketContext(ctx context.Context, handlers LobbyPacketHandlers, tag network.PacketTag, id uuid.UUID, data []byte) error {
	switch tag {
	case network.WAIT_MATCH_MAKE_PACKET_TAG:
		if handlers.MatchMake == nil {
			return ErrMissingHandler
		}
		return handlers.MatchMake.ServePacketContext(ctx, id, data)
	default:
		return ErrUnsupportedPacketTag
	}
}

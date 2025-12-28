package controller

import (
	"frontage/pkg/network"
	"github.com/google/uuid"
)

type SystemPacketHandlers struct{}

func DispatchSystemPacket(_ SystemPacketHandlers, _ network.PacketTag, _ uuid.UUID, _ []byte) error {
	return ErrUnsupportedPacketTag
}

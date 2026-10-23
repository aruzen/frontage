package controller

import (
	"context"
	"encoding/binary"
	"errors"
	"frontage/pkg/network"
	"frontage/pkg/network/repository"
	"io"
	"net"
)

var LoadHeaderFailedErr = errors.New("load header failed")

const HeaderLength = 6

func ReceiveLoop(ctx context.Context, conn net.Conn, systemChan chan network.UnsolvedPacket,
	lobbyChan chan network.UnsolvedPacket, gameChan *repository.BarrierGameChannel) error {
	defer conn.Close()
	stopClose := context.AfterFunc(ctx, func() { _ = conn.Close() })
	defer stopClose()
	header := make([]byte, HeaderLength)
	var body []byte
	for {
		if err := ctx.Err(); err != nil {
			return err
		}
		if _, err := io.ReadFull(conn, header); err != nil {
			return err
		}
		packetTag := network.PacketTag(binary.LittleEndian.Uint16(header[:2]))
		packetLength := int(binary.LittleEndian.Uint32(header[2:6]))
		if packetLength > cap(body) {
			body = make([]byte, packetLength)
		} else {
			body = body[:packetLength]
		}
		if _, err := io.ReadFull(conn, body); err != nil {
			return err
		}
		var target chan network.UnsolvedPacket
		switch packetTag >> 14 {
		case network.SystemPacketFlag:
			target = systemChan
		case network.LobbyPacketFlag:
			target = lobbyChan
		case network.GamePacketFlag:
			if gameChan != nil && gameChan.Living.Load() {
				target = gameChan.Chan
			}
		}
		if target == nil {
			continue
		}
		packet := network.UnsolvedPacket{Tag: packetTag, Body: append([]byte(nil), body...)}
		select {
		case target <- packet:
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}

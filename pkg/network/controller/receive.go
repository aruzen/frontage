package controller

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"frontage/pkg/network"
	"frontage/pkg/network/game_api"
	"log/slog"
	"net"
	"time"
)

var (
	LoadHeaderFailedErr = errors.New("load header failed")
)

const HeaderLength = 6

func ReceiveLoop(ctx context.Context, conn net.Conn, systemChan chan network.Packet, lobbyChan chan network.Packet, gameChan chan network.Packet) error {
	defer conn.Close()
	header := make([]byte, HeaderLength)
	body := make([]byte, 0)
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		loaded := 0
		for loaded < HeaderLength {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			_ = conn.SetReadDeadline(time.Now().Add(time.Second))
			read, err := conn.Read(header[loaded:])
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}
				return err
			}
			loaded += read
		}
		packetTag := binary.LittleEndian.Uint16(header[:2])
		packetLength := binary.LittleEndian.Uint32(header[2:6])
		if int(packetLength) > cap(body) {
			body = make([]byte, packetLength)
		} else {
			body = body[:packetLength]
		}
		loaded = 0
		for loaded < int(packetLength) {
			if ctx.Err() != nil {
				return ctx.Err()
			}
			_ = conn.SetReadDeadline(time.Now().Add(time.Second))
			read, err := conn.Read(body[loaded:packetLength])
			if err != nil {
				if ne, ok := err.(net.Error); ok && ne.Timeout() {
					continue
				}
				println("err body : ", read, string(body[:packetLength]))
				return err
			}
			loaded += read
		}
		var packet network.Packet
		var err error
		switch network.PacketTag(packetTag) {
		case network.ACT_EVENT_PACKET_TAG:
			tmp := game_api.ActEventPacket{}
			err = json.Unmarshal(body[:loaded], &tmp)
			packet = tmp
		case network.GAME_INITIALIZE_PACKET_TAG:
			tmp := game_api.GameInitializePacket{}
			err = json.Unmarshal(body[:loaded], &tmp)
			packet = tmp
		case network.OPPONENT_PLAYER_INITIALIZE_PACKET_TAG:
			tmp := game_api.OpponentPlayerInitializePacket{}
			err = json.Unmarshal(body[:loaded], &tmp)
			packet = tmp
		case network.MY_DECK_UPLOAD_PACKET_TAG:
			tmp := game_api.MyDeckUploadPacket{}
			err = json.Unmarshal(body[:loaded], &tmp)
			packet = tmp
		}
		if err != nil {
			slog.Error(err.Error())
			continue
		}
		if packet == nil {
			slog.Warn("unknown packet tag", "tag", packetTag)
			continue
		}
		category := packetTag >> 14
		switch network.PacketTag(category) {
		case network.SystemPacketFlag:
			systemChan <- packet
		case network.LobbyPacketFlag:
			lobbyChan <- packet
		case network.GamePacketFlag:
			gameChan <- packet
		}
	}
	return nil
}

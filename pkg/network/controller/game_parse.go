package controller

import (
	"encoding/json"
	"frontage/pkg/network"
	"frontage/pkg/network/game_api"
	"frontage/pkg/network/game_handler"
)

type GamePacketParsers struct {
	ActEvent       *game_handler.ActEventHandler
	GameInitialize *game_handler.GameInitializeHandler
	OpponentInit   *game_handler.OpponentPlayerInitializeHandler
	MyDeckUpload   *game_handler.MyDeckUploadHandler
	TurnStart      *game_handler.TurnStartHandler
	TurnEnd        *game_handler.TurnEndHandler
	TurnPass       *game_handler.TurnPassHandler
}

func ParseGamePacket(parsers GamePacketParsers, tag network.PacketTag, data []byte) (network.Packet, error) {
	switch tag {
	case network.ACT_EVENT_PACKET_TAG:
		if parsers.ActEvent == nil {
			return nil, ErrMissingHandler
		}
		packet, err := parsers.ActEvent.ServePacket(data)
		if err != nil {
			return nil, err
		}
		return packet, nil
	case network.GAME_INITIALIZE_PACKET_TAG:
		if parsers.GameInitialize == nil {
			return nil, ErrMissingHandler
		}
		size, isMySideFirst, err := parsers.GameInitialize.ServePacket(data)
		if err != nil {
			return nil, err
		}
		yourSide := 1
		if isMySideFirst {
			yourSide = 0
		}
		return game_api.GameInitializePacket{Width: size.Width, Height: size.Height, YourSide: yourSide}, nil
	case network.OPPONENT_PLAYER_INITIALIZE_PACKET_TAG:
		if parsers.OpponentInit == nil {
			return nil, ErrMissingHandler
		}
		id, name, err := parsers.OpponentInit.ServePacket(data)
		if err != nil {
			return nil, err
		}
		return game_api.OpponentPlayerInitializePacket{Id: id, Name: name}, nil
	case network.MY_DECK_UPLOAD_PACKET_TAG:
		if parsers.MyDeckUpload == nil {
			return nil, ErrMissingHandler
		}
		var packet game_api.MyDeckUploadPacket
		if err := json.Unmarshal(data, &packet); err != nil {
			return nil, err
		}
		if _, _, err := parsers.MyDeckUpload.ServePacket(data); err != nil {
			return nil, err
		}
		return packet, nil
	case network.TURN_PASS_PACKET_TAG:
		if parsers.TurnPass == nil {
			return nil, ErrMissingHandler
		}
		turn, err := parsers.TurnPass.ServePacket(data)
		if err != nil {
			return nil, err
		}
		return game_api.TurnPassPacket{Turn: turn}, nil
	case network.TURN_START_PACKET_TAG:
		if parsers.TurnStart == nil {
			return nil, ErrMissingHandler
		}
		turn, err := parsers.TurnStart.ServePacket(data)
		if err != nil {
			return nil, err
		}
		return game_api.TurnStartPacket{Turn: turn}, nil
	case network.TURN_END_PACKET_TAG:
		if parsers.TurnEnd == nil {
			return nil, ErrMissingHandler
		}
		turn, err := parsers.TurnEnd.ServePacket(data)
		if err != nil {
			return nil, err
		}
		return game_api.TurnEndPacket{Turn: turn}, nil
	default:
		return nil, ErrUnsupportedPacketTag
	}
}

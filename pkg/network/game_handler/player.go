package game_handler

import (
	"encoding/json"
	"errors"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/game_api"
	"log/slog"
)

var ErrCardNotFound = errors.New("card not found")
var ErrUnexpectedPlayer = errors.New("player is not LocalPlayer")

type ActPieceSkillHandler struct{}
type ActStructureSkillHandler struct{}
type PlayCardHandler struct{}

func NewActPieceSkillHandler() *ActPieceSkillHandler {
	return &ActPieceSkillHandler{}
}

func NewActStructureSkillHandler() *ActStructureSkillHandler {
	return &ActStructureSkillHandler{}
}

func NewPlayCardHandler() *PlayCardHandler {
	return &PlayCardHandler{}
}

func (h *ActPieceSkillHandler) ServePacket(board *model.Board, data []byte) (model.Piece, game_api.ActPieceSkillPacket, error) {
	var packet game_api.ActPieceSkillPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return nil, game_api.ActPieceSkillPacket{}, err
	}
	size := board.Size()
	pos := packet.Position
	if pos.X < 0 || pos.X >= size.Width || pos.Y < 0 || pos.Y >= size.Height {
		return nil, packet, ErrCardNotFound
	}
	piece := board.Entities()[pos.X][pos.Y]
	if piece == nil {
		return nil, packet, ErrCardNotFound
	}
	return piece, packet, nil
}

func (h *ActStructureSkillHandler) ServePacket(board *model.Board, data []byte) (model.Structure, game_api.ActStructureSkillPacket, error) {
	var packet game_api.ActStructureSkillPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return nil, game_api.ActStructureSkillPacket{}, err
	}
	size := board.Size()
	pos := packet.Position
	if pos.X < 0 || pos.X >= size.Width || pos.Y < 0 || pos.Y >= size.Height {
		return nil, packet, ErrCardNotFound
	}
	structure := board.Structures()[pos.X][pos.Y]
	if structure == nil {
		return nil, packet, ErrCardNotFound
	}
	return structure, packet, nil
}

func (h *PlayCardHandler) ServePacket(board *model.Board, data []byte) (model.Card, error) {
	var packet game_api.PlayCardPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return nil, err
	}
	players := board.Players()
	for _, p := range players {
		local, ok := p.(*model.LocalPlayer)
		if !ok {
			slog.Warn("player is not LocalPlayer", "player", p.ID())
			return nil, ErrUnexpectedPlayer
		}
		hand := local.GetDeck(model.DECK_TYPE_HAND)
		if hand != nil {
			if card, ok := hand.FindById(packet.Card.UUID); ok {
				return card, nil
			}
		}
		for _, deck := range local.Decks() {
			if deck.Deck == nil {
				continue
			}
			if card, ok := deck.Deck.FindById(packet.Card.UUID); ok {
				return card, nil
			}
		}
	}
	return nil, ErrCardNotFound
}

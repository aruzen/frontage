package game_handler

import (
	"encoding/json"
	"frontage/pkg"
	"frontage/pkg/engine/model"
	"frontage/pkg/network/game_api"
	"frontage/pkg/network/repository"
	"github.com/google/uuid"
)

type GameInitializeHandler struct {
}

type OpponentPlayerInitializeHandler struct {
}

type MyDeckUploadHandler struct {
	cardRepo *repository.CardRepository
}

func NewInitializeHandler() *GameInitializeHandler {
	return &GameInitializeHandler{}
}

func (handler *GameInitializeHandler) ServePacket(data []byte) (pkg.Size, bool /*isMySideFirst*/, error) {
	var packet game_api.GameInitializePacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return pkg.Size{}, false, err
	}

	boardSize := pkg.Size{Width: packet.Width, Height: packet.Height}
	isMySideFirst := packet.YourSide == 0
	return boardSize, isMySideFirst, nil
}

func NewOpponentPlayerInitializeHandler() *OpponentPlayerInitializeHandler {
	return &OpponentPlayerInitializeHandler{}
}

func (handler *OpponentPlayerInitializeHandler) ServePacket(data []byte) (uuid.UUID, string, error) {
	var packet game_api.OpponentPlayerInitializePacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return uuid.Nil, "", err
	}
	return packet.Id, packet.Name, nil
}

func NewMyDeckUploadHandler(cardRepo *repository.CardRepository) MyDeckUploadHandler {
	return MyDeckUploadHandler{
		cardRepo: cardRepo,
	}
}

func (h *MyDeckUploadHandler) ServePacket(data []byte) (model.Cards, model.Cards, error) {
	var packet game_api.MyDeckUploadPacket
	if err := json.Unmarshal(data, &packet); err != nil {
		return nil, nil, err
	}
	instantiateDeck := func(tags []string) (model.Cards, error) {
		deck := make(model.Cards, len(tags))
		for i, tag := range tags {
			cardTemplate, err := h.cardRepo.Find(pkg.ItemTag(tag))
			if err != nil {
				return nil, err
			}
			deck[i] = cardTemplate.CardMirror(uuid.New())
		}
		return deck, nil
	}

	mainDeck, err := instantiateDeck(packet.MainDeck)
	if err != nil {
		return nil, nil, err
	}
	subDeck, err := instantiateDeck(packet.SubDeck)
	if err != nil {
		return nil, nil, err
	}
	return mainDeck, subDeck, nil
}

package handler

import (
	"errors"
	"frontage/pkg/data"
	"frontage/pkg/engine/model"
)

func InstantiationCard(c data.Card) (model.Card, error) {
	card, ok := data.CardTable[c.Name]
	if !ok {
		return nil, errors.New("Card not found")
	}
	mirror := card.Mirror(c.UUID)
	return mirror, nil
}

func InstantiationDeck(cardsData []data.Card) model.Cards {
	deck := make(model.Cards, len(cardsData))
	for i, cardData := range cardsData {
		card, err := InstantiationCard(cardData)
		if err != nil {
			panic(err)
		}
		deck[i] = card
	}
	return deck
}

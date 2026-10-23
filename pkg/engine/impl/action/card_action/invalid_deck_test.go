package card_action

import (
	"frontage/pkg/engine/model"
	"testing"
)

func TestPieceActionStateRejectsInvalidDeckType(t *testing.T) {
	board, player, piece := newTestBoardWithCard(t)
	for _, tc := range []struct {
		name     string
		deckType model.DeckType
	}{
		{"negative", model.DeckType(-1)},
		{"out of range", model.DeckType(999)},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var state PieceActionState
			err := state.FromMap(board, map[string]interface{}{
				"holder_id": player.ID().String(),
				"card_id":   piece.Id().String(),
				"deck_type": int(tc.deckType),
				"value":     1,
			})
			if err == nil {
				t.Fatal("invalid deck type was accepted")
			}
		})
	}
}

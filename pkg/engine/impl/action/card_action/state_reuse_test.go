package card_action

import (
	"testing"

	"frontage/pkg"
	"frontage/pkg/engine/impl/card"
	"frontage/pkg/engine/logic"
	"frontage/pkg/engine/model"
	"github.com/google/uuid"
)

func TestCardStateReuseClearsReferences(t *testing.T) {
	board, player, piece := newTestBoardWithCard(t)
	operate := &PieceActionState{holderID: player.ID(), cardID: piece.Id(), deckType: model.DECK_TYPE_MAIN, value: 1}
	summon := NewPieceSummonActionState(player, model.DECK_TYPE_MAIN, piece, pkg.Point{X: 1, Y: 1})
	for _, tc := range []struct {
		name   string
		state  logic.ActionState
		holder func() *model.LocalPlayer
		piece  func() card.Piece
	}{
		{"operate", operate, func() *model.LocalPlayer { return operate.holder }, func() card.Piece { return operate.piece }},
		{"summon", &summon, func() *model.LocalPlayer { return summon.holder }, func() card.Piece { return summon.piece }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			valid := tc.state.ToMap()
			for _, scenario := range []string{"nil board", "missing player", "missing card", "invalid input"} {
				t.Run(scenario, func(t *testing.T) {
					if err := tc.state.FromMap(board, valid); err != nil {
						t.Fatal(err)
					}
					if tc.holder() != player || tc.piece() != piece {
						t.Fatal("initial references were not restored")
					}
					input := tc.state.ToMap()
					nextBoard := board
					switch scenario {
					case "nil board":
						nextBoard = nil
					case "missing player":
						input["holder_id"] = uuid.New().String()
					case "missing card":
						input["card_id"] = uuid.New().String()
					case "invalid input":
						input["holder_id"] = "invalid"
					}
					err := tc.state.FromMap(nextBoard, input)
					if (err != nil) != (scenario != "nil board") {
						t.Errorf("unexpected restore error: %v", err)
					}
					if tc.piece() != nil {
						t.Error("card reference from the previous restore survived")
					}
					if scenario == "missing card" {
						if tc.holder() != player {
							t.Error("current holder was not resolved")
						}
					} else if tc.holder() != nil {
						t.Error("holder reference from the previous restore survived")
					}
				})
			}
		})
	}
}

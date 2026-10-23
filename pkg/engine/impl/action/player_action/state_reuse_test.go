package player_action

import (
	"testing"

	"github.com/google/uuid"
)

func TestDrawStateReuseClearsReferences(t *testing.T) {
	board, player, _, _, _ := newTestBoardWithLocalPlayer(t)
	for _, scenario := range []string{"nil board", "missing player", "invalid input"} {
		t.Run(scenario, func(t *testing.T) {
			state := NewDrawActionState(player)
			input := state.ToMap()
			nextBoard := board
			switch scenario {
			case "nil board":
				nextBoard = nil
			case "missing player":
				input["player_id"] = uuid.New().String()
			case "invalid input":
				input["player_id"] = "invalid"
			}
			err := state.FromMap(nextBoard, input)
			if (err != nil) != (scenario == "invalid input") {
				t.Errorf("unexpected restore error: %v", err)
			}
			if state.Player() != nil {
				t.Fatal("player reference from the previous restore survived")
			}
		})
	}
}
